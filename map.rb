require 'sinatra'
require 'sinatra/json'
require 'exifr/jpeg'
require 'csv'
require 'plus_codes/open_location_code'
require 'json'

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
  send_file File.join(File.absolute_path('assets/photos/') + '/' + file)
end

def photo_gps_list(region)
  table = CSV.open("assets/graffiti.csv", headers: true).map(&:to_h)
  olc = PlusCodes::OpenLocationCode.new
  photos_url = '/assets/photos/'
  all = region == "all"
  result = table.map do |photo|
    plus_code = olc.encode(Float(photo["latitude"]), Float(photo["longitude"]), 16)
    {
      type: 'Feature',
      "geometry": { "type": "Point", "coordinates": [photo["longitude"], photo["latitude"]]},
      "properties": {
        "image": File.join(photos_url, photo["filename"]),
        "ipfs": photo["ipfs"],
        "surface": photo["surface"],
        "url": File.join(photos_url, photo["filename"]),
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
  all = region == "all"
  starts_with_region = "^#{region}"
  olc = PlusCodes::OpenLocationCode.new
  table = CSV.open("assets/graffiti.csv", headers: true).map(&:to_h)
  table.map { |photo|
    latitude = Float(photo["latitude"])
    longitude = Float(photo["longitude"])
    photo["plus_code"] = olc.encode(latitude, longitude, 16)
    photo["latitude"] = latitude
    photo["longitude"] = longitude
    photo["local"] = "/assets/photos/#{photo["filename"]}"
    photo["url"] = "https://ipfs.io/ipfs/#{photo["ipfs"]}"
  }
  table.map { |record| record if all or record["plus_code"].match(starts_with_region) }.compact
end
