Requirements:

- Ruby (https://www.ruby-lang.org/en/)
- IPFS (https://ipfs.io/)

Run:

```sh
git clone https://github.com/romanblanco/graffiti.git
cd graffiti/
mkdir assets
cd assets
ipfs get QmeNNGcqg12BWoyHWJ1Aa6WaeTrct5WHjPpQ1LUGip7se1
ln -s QmeNNGcqg12BWoyHWJ1Aa6WaeTrct5WHjPpQ1LUGip7se1 graffiti.csv
ipfs get QmdWeEuqA6gHACFGYd8yfiwyX8QGrQ7GzxRDdQPxf3VZxA
ln -s QmdWeEuqA6gHACFGYd8yfiwyX8QGrQ7GzxRDdQPxf3VZxA/ photos
gem install bundler
bundle install
bundle exec ruby map.rb -o 0.0.0.0
```

`http://0.0.0.0:4567/`

![index](/index.png "index")

`http://0.0.0.0:4567/detail/8FXR5JW`

![detail](/detail.png "detail")

`http://0.0.0.0:4567/api/8FXR5JW`

![api](/api.png "api")
