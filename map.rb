require 'sinatra'
require 'sinatra/json'
require 'ipfs/client'
require 'plus_codes/open_location_code'
require 'json'
require 'csv'

set :token, File.read('./TOKEN').strip
set :static, true
set :root, File.dirname(__FILE__)
set :public_folder, Proc.new { File.join(root, "assets/photos/") }

class Resource
  attr_accessor :ipfs
  attr_accessor :graffiti_data
  attr_accessor :images_data

  def initialize
    @ipfs = IPFS::Client.default # uses localhost and port 5001
    @images_nodes = [
      'QmdWeEuqA6gHACFGYd8yfiwyX8QGrQ7GzxRDdQPxf3VZxA',
    ]
    @graffiti_files = [
      'QmeNNGcqg12BWoyHWJ1Aa6WaeTrct5WHjPpQ1LUGip7se1',
    ]
    @graffiti_data = []
    @images_data = []
    load_resource
  end

  private

  def load_resource
    # load information about images available from ipfs nodes
    @images_nodes.each do |node|
      @images_data.push(
        @ipfs.ls(node).map do |node|
          node.links.map { |link| { ipfs: link.hashcode, size: link.size } }
        end.first)
    end
    # load data describing graffiti images
    @graffiti_files.each do |file|
      @graffiti_data.push(CSV.parse(@ipfs.cat(file), headers: true).map(&:to_h)).first
    end
  end
end

@@resource = Resource.new

get '/' do
  erb :index, :locals => {
    :data => markers('all'),
    :token => settings.token
  }
end

get '/detail/?:region?' do |r|
  erb :detail, :locals => {
    :data => markers(params['region']),
    :token => settings.token
  }
end

get '/api/?:region?' do |r|
  json(api(params['region']))
end

get "/assets/photos/:file" do |file|
  filename = file + '.jpg'
  file_path = 'assets/photos/'
  File.write(file_path + filename, @@resource.ipfs.cat(file))
  send_file File.join(File.absolute_path(file_path) + '/' + filename)
end

def api(search)
  olc = PlusCodes::OpenLocationCode.new
  images = @@resource.images_data.flatten
  descriptions = @@resource.graffiti_data.flatten
  all = search == "all"
  starts_with_region = "^#{search}".upcase
  result = images.map do |image|
    image_details = descriptions.select do |description|
      image[:ipfs] == description['ipfs']
    end
    image[:url] = "https://ipfs.io/ipfs/#{image[:ipfs]}"
    if !image_details.empty?
      data = image_details.first
      image[:date] = data['date']
      image[:latitude] = data['latitude']
      image[:longitude] = data['longitude']
      image[:plus_code] = olc.encode(
        Float(data['latitude']),
        Float(data['longitude']),
        16)
      image[:surface] = data["surface"]
    else
      image[:date] = ''
      image[:latitude] = ''
      image[:longitude] = ''
      image[:plus_code] = ''
      image[:surface] = ''
    end
    image if all or image[:plus_code].match(starts_with_region)
  end.compact.sort_by { |image| image[:plus_code] }
  result
end

def markers(region)
  photos_url = '/assets/photos/'
  api(region).map do |photo|
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
        "plus": photo[:plus_code],
        "marker-symbol": "art-gallery",
        "marker-color": photo[:surface] != nil ? "#00FF00" : "#000000",
        "marker-size": "medium",
      }
    } if photo[:plus_code] != ''
  end.compact.to_json
end
