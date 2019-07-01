Requirements:

- Ruby (https://www.ruby-lang.org/en/)
- IPFS (https://ipfs.io/)


Resources:

- Geotagged photos:
  - https://explore.ipld.io/#/explore/QmdWeEuqA6gHACFGYd8yfiwyX8QGrQ7GzxRDdQPxf3VZxA
- CSV file with graffiti description
  - https://explore.ipld.io/#/explore/QmeNNGcqg12BWoyHWJ1Aa6WaeTrct5WHjPpQ1LUGip7se1

Run:

```sh
ipfs pin add QmdWeEuqA6gHACFGYd8yfiwyX8QGrQ7GzxRDdQPxf3VZxA
ipfs pin add QmeNNGcqg12BWoyHWJ1Aa6WaeTrct5WHjPpQ1LUGip7se1
git clone https://github.com/romanblanco/graffiti.git
mkdir graffiti/assets/photos
cd graffiti/
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
