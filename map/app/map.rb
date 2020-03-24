require 'sinatra'
require 'sinatra/json'
require 'ipfs-http-client'
require 'plus_codes/open_location_code'
require 'json'
require 'exifr/jpeg'
require 'csv'
require 'geocoder'

set :token, File.read('./TOKEN').strip
set :static, true
set :root, File.dirname(__FILE__)
set :public_folder, Proc.new { File.join(root, "assets/photos/") }

sleep(2) # to make sure IPFS is started by that time

class Resource
  attr_accessor :ipfs
  attr_accessor :graffiti_data
  attr_accessor :images_data

  def initialize
    @ipfs = Ipfs::Client.new 'http://ipfs:5001'
    @images_nodes = [
      'QmdWeEuqA6gHACFGYd8yfiwyX8QGrQ7GzxRDdQPxf3VZxA',
      'QmYVGFdAxxXYK2E8Ub8Xoe69YgAx19utAQZ639noYCvNxU',
    ]
    @graffiti_files = [
      #@ipfs.cat('QmeNNGcqg12BWoyHWJ1Aa6WaeTrct5WHjPpQ1LUGip7se1'),
      File.read('graffiti.json')
    ]
    @offline_area = [
      '8FXR5JW',
    ]
    @graffiti_data = []
    @images_data = []
    load_resource
  end

  private

  def load_resource
    images_data = []
    graffiti_data = []

    # load information about images available from ipfs nodes
    @images_nodes.each do |node|
      p node
      images_data.push(
        @ipfs.ls(node).map do |node|
          node.links.map { |link| { ipfs: link.hashcode, size: link.size } }
        end.first)
      @images_data = images_data.flatten.uniq
    end

    # load data describing graffiti images
    @graffiti_files.each do |file|
      graffiti_data.push(JSON.parse(file))
      @graffiti_data = graffiti_data.flatten
    end
    @graffiti_data.map! do |hash|
      hash.inject({}) { |memo, (k, v)| memo[k.to_sym] = v; memo }
    end

    # pre-download images
    @graffiti_data.each do |image|
      if !image[:olc].empty? &&
         @offline_area.any? do |area|
           image[:olc].match("^#{area.upcase}")
         end
          image[:olc].match("^#{@offline_area[0]}".upcase)
        unless File.file?("assets/photos/#{image[:ipfs]}.jpg")
          p "downloading #{image[:ipfs]}"
          File.write("assets/photos/#{image[:ipfs]}.jpg", @ipfs.cat(image[:ipfs]))
        end
      end
    end

    undescribed = @images_data.map do |undescribed|
      undescribed if @graffiti_data.select do |image|
        image[:ipfs] == undescribed[:ipfs]
      end.empty?
    end.compact

    undescribed.each { |image|
      p "graffiti.json should be updated, missing #{image[:ipfs]}"
      File.write("assets/photos/#{image[:ipfs]}.jpg", @ipfs.cat(image[:ipfs]))
    }

    ipfs_download = undescribed.map { |image|
      exif = EXIFR::JPEG.new("assets/photos/#{image[:ipfs]}.jpg")
      { :ipfs => image[:ipfs],
        :date =>  exif.date_time.to_s,
        :filename => nil,
        :longitude => exif&.gps&.longitude.to_s,
        :latitude => exif&.gps&.latitude.to_s }
    }
    ipfs_download.each { |img| @graffiti_data.push(img) }
  end
end

p Dir.pwd
@@resource = Resource.new

get '/' do
  p Dir.pwd
  erb :cluster, :locals => {
    :data => markers(api(:type => :all)),
    :token => settings.token
  }
end

get '/detail/:region' do |r|
  erb :detail, :locals => {
    :data => markers(api(:type => :city,
                         :request => Geocoder.search(params['region']).first.data["boundingbox"])),
    :token => settings.token
  }
end

get '/city/?:city?' do |r|
  erb :cluster, :locals => {
    :data => markers(api(
      :type => :city,
      :request => Geocoder.search(params['city']).first.data["boundingbox"])),
    :token => settings.token
  }
end

get '/tag/?:tag?' do |r|
  erb :cluster, :locals => {
    :data => markers(api(
      :type => :tag,
      :request => params['tag'].split(','))),
    :token => settings.token
  }
end

get '/api/?:ipfs?' do |r|
  if params.empty?
    json(api(:type => :all))
  else
    json(api(:type => :ipfs, :request => params['ipfs']))
  end
end

get '/api/city/?:city?' do |r|
  json(api(
    :type => :city,
    :request => Geocoder.search(params['city']).first.data["boundingbox"]))
end

get '/api/tag/?:tag?' do |r|
  json(api(
    :type => :tag,
    :request => params['tag'].split(',')))
end

get '/api/plus/?:region?' do |r|
  json(api(:type => :plus, :request => params['region']))
end

get "/assets/photos/:file" do |file|
  filename = file + '.jpg'
  file_path = 'assets/photos/'
  File.write(file_path + filename, @@resource.ipfs.cat(file))
  send_file File.join(File.absolute_path(file_path) + '/' + filename)
end

def api(search_params)
  olc = PlusCodes::OpenLocationCode.new
  result = @@resource.images_data.map do |image|
    image_details = @@resource.graffiti_data.select do |description|
      image[:ipfs] == description[:ipfs]
    end
    if !image_details.empty? && !(image_details.first[:latitude] == "" ||
                                  image_details.first[:longitude] == "")
      data = image_details.first
      image[:date] = data[:date]
      image[:latitude] = data[:latitude]
      image[:longitude] = data[:longitude]
      image[:olc] = olc.encode(
        Float(data[:latitude]),
        Float(data[:longitude]),
        16)
      image[:surface] = data[:surface]
      image[:tag] = !data.nil? && !data[:tag].nil? ? data[:tag]  : []
    else
      image[:date] = ''
      image[:latitude] = ''
      image[:longitude] = ''
      image[:olc] = ''
      image[:surface] = nil
      image[:tag] = !data.nil? && !data[:tag].nil? ? data[:tag]  : []
    end
    image if search(image, search_params)
  end.compact.sort_by { |image| image[:olc] }
  result
end

def search(image, params)
  case params[:type]
  when :all
    true
  when :ipfs
    image[:ipfs] == params[:request]
  when :plus
    image[:olc].match("^#{params[:request]}".upcase)
  when :city
    city(image[:latitude], image[:longitude], params[:request])
  when :tag
    p image[:tag]

    !((image[:tag] & params[:request]).empty?) if !image[:tag].nil?
  end
end

def city(img_lat, img_lon, city)
  min_lat, max_lat, min_lon, max_lon = city.map(&:to_f)
  img_lat = img_lat.to_f
  img_lon = img_lon.to_f
  lat_ok = (min_lat <= img_lat && img_lat <= max_lat)
  lon_ok = (min_lon <= img_lon && img_lon <= max_lon)
  lat_ok && lon_ok
end

def markers(data)
  photos_url = '/assets/photos/'
  data.map do |photo|
    {
      type: 'Feature',
      "geometry": { "type": "Point", "coordinates": [photo[:longitude], photo[:latitude]]},
      "properties": {
        "image": File.join(photos_url, photo[:ipfs]),
        "ipfs": photo[:ipfs],
        "surface": photo[:surface],
        "url": File.join(photos_url, photo[:ipfs]),
        "date": photo[:date],
        "gps_longitude": photo[:longitude],
        "gps_latitude": photo[:latitude],
        "tag": photo[:tag].join(','),
        "olc": photo[:olc],
        "marker-symbol": "art-gallery",
        "marker-color": photo[:surface] != nil ? "#0088ce" : "#000000",
        "marker-size": "medium",
      }
    } if photo[:olc] != ''
  end.compact.to_json
end