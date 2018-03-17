require 'sinatra'
require 'exifr/jpeg'
require 'json'

require_relative 'app.rb'

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

get "/assets/photos/:file" do |file|
  send_file File.join(File.absolute_path('assets/photos/') + '/' + file)
end

def photo_gps_list
  photos_dir = File.absolute_path('assets/photos/')
  photos_url = '/assets/photos/'
  photos = Dir.entries(photos_dir).select { |filename| filename if filename.end_with?('jpg') }
  result = photos.map do |photo|
    photo_exif = EXIFR::JPEG.new(photos_dir + '/' + photo)
    {
      type: 'Feature',
      "geometry": { "type": "Point", "coordinates": [photo_exif.gps.longitude, photo_exif.gps.latitude]},
      "properties": {
        "image": File.join(photos_url, photo),
        "url": File.join(photos_url, photo),
        "date": photo_exif.exif.date_time_original,
        "marker-symbol": "art-gallery",
        "marker-color": "#000000",
        "marker-size": "medium",
      }
    } unless photo_exif.gps.nil?
  end
  result.compact.to_json
end
