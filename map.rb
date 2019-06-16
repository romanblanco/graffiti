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
    :data => photo_gps_list,
    :token => settings.token
  }
end

get '/detail' do
  erb :detail, :locals => {
    :data => photo_gps_list,
    :token => settings.token
  }
end

get '/export/:region' do |r|
  json(export(params['region']))
end


get "/assets/photos/:file" do |file|
  send_file File.join(File.absolute_path('assets/photos/') + '/' + file)
end

def photo_gps_list
  table = CSV.open("assets/graffiti.csv", headers: true).map(&:to_h)
  olc = PlusCodes::OpenLocationCode.new
  photos_url = '/assets/photos/'
  result = table.map do |photo|
    {
      type: 'Feature',
      "geometry": { "type": "Point", "coordinates": [photo["longitude"], photo["latitude"]]},
      "properties": {
        "image": File.join(photos_url, photo["filename"]),
        "ipfs": photo["ipfs"],
        "url": File.join(photos_url, photo["filename"]),
        "date": photo["date"],
        "gps_longitude": photo["longitude"],
        "gps_latitude": photo["latitude"],
        "plus": olc.encode(Float(photo["latitude"]), Float(photo["longitude"]), 16),
        "marker-symbol": "art-gallery",
        "marker-color": "#000000",
        "marker-size": "medium",
      }
    }
  end
  result.compact.to_json
end

def export(region)
  all = region == "all"
  starts_with_region = "^#{region}"
  olc = PlusCodes::OpenLocationCode.new
  table = CSV.open("assets/graffiti.csv", headers: true).map(&:to_h)
  table.map { |photo|
    photo["plus_code"] = olc.encode(Float(photo["latitude"]), Float(photo["longitude"]), 16)
    photo["url"] = "https://ipfs.io/ipfs/#{photo["ipfs"]}"
  }
  table.map { |record| record if all or record["plus_code"].match(starts_with_region) }.compact
end
