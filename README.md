Requirements:

- Ruby (https://www.ruby-lang.org/en/)
- IPFS (https://ipfs.io/)

Run:

```sh
git clone https://github.com/romanblanco/EXIF-GPS-map.git
cd EXIF-GPS-map/
ipfs get QmeNNGcqg12BWoyHWJ1Aa6WaeTrct5WHjPpQ1LUGip7se1
ln -s QmeNNGcqg12BWoyHWJ1Aa6WaeTrct5WHjPpQ1LUGip7se1 assets/graffiti.csv
ipfs get QmdWeEuqA6gHACFGYd8yfiwyX8QGrQ7GzxRDdQPxf3VZxA
ln -s QmdWeEuqA6gHACFGYd8yfiwyX8QGrQ7GzxRDdQPxf3VZxA/ assets/photos
gem install bundler
bundle install
bundle exec ruby map.rb -o 0.0.0.0
```
