require 'sinatra'
require 'sinatra/json'
require 'exifr/jpeg'
require 'csv'
require 'plus_codes/open_location_code'
require 'json'
require 'ipfs/client'

set :token, File.read('./TOKEN').strip
set :static, true
set :root, File.dirname(__FILE__)
set :public_folder, Proc.new { File.join(root, "assets/photos/") }

get '/' do
  erb :index, :locals => {
    :data => photo_gps_list('all'),
    :token => settings.token
  }
end

get '/detail/?:region?' do |r|
  erb :detail, :locals => {
    :data => photo_gps_list(params['region']),
    :token => settings.token
  }
end

get '/api/?:region?' do |r|
  json(export(params['region']))
end

get "/assets/photos/:file" do |file|
  filename = file + '.jpg'
  file_path = 'assets/photos/'
  client = IPFS::Client.default # => uses localhost and port 5001
  File.write(file_path + filename, client.cat(file))
  send_file File.join(File.absolute_path(file_path) + '/' + filename)
end

get '/download' do
  erb :download, :locals => {
    :ipfs => download
  }
end

post "/getFromIpfs" do
  file = params['ipfs_content']
  filename = file + '.jpg'
  file_path = 'assets/photos/'
  client = IPFS::Client.default # => uses localhost and port 5001
  File.write(file_path + filename, client.cat(file))
  redirect '/download'
end

def download
  client = IPFS::Client.default # => uses localhost and port 5001
  content = client.ls 'QmdWeEuqA6gHACFGYd8yfiwyX8QGrQ7GzxRDdQPxf3VZxA'
  ipfs_content = content.map { |node| node.links.map { |link| {ipfs: link.hashcode, size: link.size} } }.first
  table = CSV.parse(client.cat('QmeNNGcqg12BWoyHWJ1Aa6WaeTrct5WHjPpQ1LUGip7se1'), headers: true).map(&:to_h)

  result = ipfs_content.map { |ipfs|
    tbl = table.select { |tbl| tbl["ipfs"] == ipfs[:ipfs] }
    if !tbl.empty?
      data = tbl.first
      ipfs[:date] = data["date"]
      ipfs[:latitude] = data["latitude"]
      ipfs[:longitude] = data["longitude"]
      ipfs[:surface] = data["surface"]
    else
      ipfs[:date] = ""
      ipfs[:latitude] = ""
      ipfs[:longitude] = ""
      ipfs[:surface] = ""
    end
    ipfs
  }
  result
end

def photo_gps_list(region)
  olc = PlusCodes::OpenLocationCode.new
  client = IPFS::Client.default # => uses localhost and port 5001
  table = CSV.parse(client.cat('QmeNNGcqg12BWoyHWJ1Aa6WaeTrct5WHjPpQ1LUGip7se1'), headers: true).map(&:to_h)
  photos_url = '/assets/photos/'
  all = region == "all"
  result = table.map do |photo|
    plus_code = olc.encode(Float(photo["latitude"]), Float(photo["longitude"]), 16)
    {
      type: 'Feature',
      "geometry": { "type": "Point", "coordinates": [photo["longitude"], photo["latitude"]]},
      "properties": {
        "image": File.join(photos_url, photo["ipfs"]),
        "ipfs": photo["ipfs"],
        "surface": photo["surface"],
        "url": File.join(photos_url, photo["ipfs"]),
        "date": photo["date"],
        "gps_longitude": photo["longitude"],
        "gps_latitude": photo["latitude"],
        "plus": plus_code,
        "marker-symbol": "art-gallery",
        "marker-color": photo["surface"] != nil ? "#00FF00" : "#000000",
        "marker-size": "medium",
      }
    } if all or plus_code.match("^#{region}")
  end
  result.compact.to_json
end

def export(region)
  client = IPFS::Client.default # => uses localhost and port 5001
  all = region == "all"
  starts_with_region = "^#{region}"
  olc = PlusCodes::OpenLocationCode.new
  table = CSV.parse(client.cat('QmeNNGcqg12BWoyHWJ1Aa6WaeTrct5WHjPpQ1LUGip7se1'), headers: true).map(&:to_h)
  table.map { |photo|
    latitude = Float(photo["latitude"])
    longitude = Float(photo["longitude"])
    photo["plus_code"] = olc.encode(latitude, longitude, 16)
    photo["latitude"] = latitude
    photo["longitude"] = longitude
    photo["url"] = "https://ipfs.io/ipfs/#{photo["ipfs"]}"
  }
  table.map { |record| record if all or record["plus_code"].match(starts_with_region) }.compact
end
