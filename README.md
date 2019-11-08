http://185.8.166.79/

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
ipfs get QmdWeEuqA6gHACFGYd8yfiwyX8QGrQ7GzxRDdQPxf3VZxA
ipfs get QmeNNGcqg12BWoyHWJ1Aa6WaeTrct5WHjPpQ1LUGip7se1
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

Sharing directory on IPFS

https://stackoverflow.com/questions/39803954/ipfs-how-to-add-a-file-to-an-existing-folder/43902411#43902411

`$ ipfs add -r ./graffiti`

```
...
added QmWbSC8zuUPHrL3aGLcg4L3rdZeN6x8fcvnDVgsBrYHciH graffiti/IMG_20190301_200805.jpg
added QmSKGZjpvKT44H2S8Zp9UkZD2MieUGnVMxjNLhvUUP6kdN graffiti/IMG_20190301_200810.jpg
added QmUne4d3VYNrSrEM1HWiJadzwvVYAmgxHaupcX52ZUJozn graffiti/IMG_20190702_134904.jpg
added QmPn6mSJueqHgaQipX4MKoUwnUv23whANoJ4z8vP7XhMpQ graffiti/IMG_20190702_135027.jpg
added QmRSnqHbhXyR9puTDrXcwT6QUN9Vhgpmusb46Fb7Pdrq61 graffiti/IMG_20190702_162303.jpg
added QmSYmRpgujRXtXjFwQo8wAhNDKDWfk9QxLw6E5BHrUhvyj graffiti/IMG_20190703_100410.jpg
added QmVTEdBPCSWvjqjozPCqX9TXa5spRpcnPEnNewPDMSQ3bE graffiti/IMG_20190703_100550.jpg
added QmPSwPh1oKs9Vaav8epnnFcNMt2yRyprxTYyaSTJeYHYVC graffiti/IMG_20190703_100605.jpg
added QmdEDX3E8BUNkprx4npoZLu5Ty5S9yELKhy1bZkqMo5zpd graffiti/IMG_20190703_100959.jpg
added QmViCAsm52MyMCwgkXt1TPiNuFrbbKFATMd3rV6KrM5uHg graffiti/IMG_20190703_101017.jpg
added QmZFquTu9AZ6qyJD3QcknwrK2sg3DBV5zZKidRAtUHi1sS graffiti/IMG_20190703_101030.jpg
added QmYVGFdAxxXYK2E8Ub8Xoe69YgAx19utAQZ639noYCvNxU graffiti
 2.40 GiB / 2.40 GiB [================================================================] 100.00%
 ```
`$ ipfs ls`

```
QmWbSC8zuUPHrL3aGLcg4L3rdZeN6x8fcvnDVgsBrYHciH 382667   IMG_20190301_200805.jpg
QmSKGZjpvKT44H2S8Zp9UkZD2MieUGnVMxjNLhvUUP6kdN 346472   IMG_20190301_200810.jpg
QmUne4d3VYNrSrEM1HWiJadzwvVYAmgxHaupcX52ZUJozn 3633530  IMG_20190702_134904.jpg
QmPn6mSJueqHgaQipX4MKoUwnUv23whANoJ4z8vP7XhMpQ 3847445  IMG_20190702_135027.jpg
QmRSnqHbhXyR9puTDrXcwT6QUN9Vhgpmusb46Fb7Pdrq61 3350850  IMG_20190702_162303.jpg
QmSYmRpgujRXtXjFwQo8wAhNDKDWfk9QxLw6E5BHrUhvyj 3180312  IMG_20190703_100410.jpg
QmVTEdBPCSWvjqjozPCqX9TXa5spRpcnPEnNewPDMSQ3bE 4551676  IMG_20190703_100550.jpg
QmPSwPh1oKs9Vaav8epnnFcNMt2yRyprxTYyaSTJeYHYVC 4668407  IMG_20190703_100605.jpg
QmdEDX3E8BUNkprx4npoZLu5Ty5S9yELKhy1bZkqMo5zpd 2917300  IMG_20190703_100959.jpg
QmViCAsm52MyMCwgkXt1TPiNuFrbbKFATMd3rV6KrM5uHg 2988945  IMG_20190703_101017.jpg
QmZFquTu9AZ6qyJD3QcknwrK2sg3DBV5zZKidRAtUHi1sS 3643637  IMG_20190703_101030.jpg
```

`$ ipfs pin add QmYVGFdAxxXYK2E8Ub8Xoe69YgAx19utAQZ639noYCvNxU`

`$ ipfs name publish /ipfs/QmYVGFdAxxXYK2E8Ub8Xoe69YgAx19utAQZ639noYCvNxU`

 ```
 Published to QmZCtha6AHsaNRV1LScMXLewjx9M2imPTxKaP2ty2TJ219: /ipfs/QmYVGFdAxxXYK2E8Ub8Xoe69YgAx19utAQZ639noYCvNxU
 ```

