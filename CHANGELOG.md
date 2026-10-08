# Changelog

## [1.14.0](https://github.com/stfsy/go-api-kit/compare/v1.13.1...v1.14.0) (2026-10-08)


### Features

* implement performance improvements ([e0c8688](https://github.com/stfsy/go-api-kit/commit/e0c8688df0af3ef01e85dab1ceabfa748b4a5fdc))
* **server:** populate CrossOriginProtection from CORS origins and add context timeout handler ([7a5a8a2](https://github.com/stfsy/go-api-kit/commit/7a5a8a23fa39cd2b0864713d92e7e7661b8af2c0))


### Bug Fixes

* **config:** ensure thread-safe lazy loading using sync.Once ([65244be](https://github.com/stfsy/go-api-kit/commit/65244be8650fb31f7a9a788f56dba8e8b2cee34f))
* **handlers:** marshal in memory before writing and allow empty payloads in validating handler ([9829182](https://github.com/stfsy/go-api-kit/commit/9829182bf801e29627445b2a61ab40dc68ff9f99))
* **middlewares:** allow Content-Length 0 for write requests per RFC 9110 ([8ddcdc7](https://github.com/stfsy/go-api-kit/commit/8ddcdc76c3a9a1f5c6dba234e0eae897473be612))
* **security:** set X-XSS-Protection to 0 and add Content-Security-Policy header ([179b6b6](https://github.com/stfsy/go-api-kit/commit/179b6b6e553fcf27183d9cb535115238a6153ccc))
* **validation:** limit struct recursion depth and return defensive copies from cache ([0327225](https://github.com/stfsy/go-api-kit/commit/0327225e243ea649280e374f9bdca35b96ec149d))


### Performance Improvements

* **config:** return configuration pointer to avoid pass-by-value struct copies ([5da4cfb](https://github.com/stfsy/go-api-kit/commit/5da4cfb5686319e2d2edd4d719d14b1c980a46a9))
* **handlers:** optimize validation allocations and pre-render static error payloads ([c15b0a8](https://github.com/stfsy/go-api-kit/commit/c15b0a8719d3e31d1647d06f18033a0ea259f92d))
* **middlewares:** precompute header slices and optimize content-type checking ([74c7ecc](https://github.com/stfsy/go-api-kit/commit/74c7ecc6fc4aceb12584ed88f89b74f2d5c5ffea))
* **server:** standardize address bindings and remove access log from default stack ([df64c7b](https://github.com/stfsy/go-api-kit/commit/df64c7bfe7bc8b0298601fc75ff3269284624749))
* **utils:** implement lock-free reads in LimitedCache using atomic.Pointer ([8f763d6](https://github.com/stfsy/go-api-kit/commit/8f763d6626cb8d44fa23e96e6646352661abfddd))

## [1.13.1](https://github.com/stfsy/go-api-kit/compare/v1.13.0...v1.13.1) (2026-07-28)


### Bug Fixes

* **validation:** use reflect.Pointer to satisfy govet inline check ([169184b](https://github.com/stfsy/go-api-kit/commit/169184ba90bf288d630755cab11b442c6de27f3b))

## [1.13.0](https://github.com/stfsy/go-api-kit/compare/v1.12.0...v1.13.0) (2025-12-05)


### Features

* add csrf protection ([3268add](https://github.com/stfsy/go-api-kit/commit/3268add9e9965e15ba3a8f8ea57ba8797999f0ec))
* enable csrf protection ([d319ab5](https://github.com/stfsy/go-api-kit/commit/d319ab59c672f65efe372a02914b9a1c2a8caee6))
* use latest negroni release ([3f868fb](https://github.com/stfsy/go-api-kit/commit/3f868fb57d011013d37a9ced1980d3c4b0b42ca4))

## [1.12.0](https://github.com/stfsy/go-api-kit/compare/v1.11.0...v1.12.0) (2025-10-11)


### Features

* use mime.ParseMediaType to parse content type ([327c236](https://github.com/stfsy/go-api-kit/commit/327c236e498f2d057837eb8c1a8aa61f6ce58008))


### Bug Fixes

* fix body detection for http delete ([bad05bb](https://github.com/stfsy/go-api-kit/commit/bad05bb0c80c9a9c54da890648b6804660a020f2))
* if content-length is -1 for delete then check for transfer encoding header ([d120255](https://github.com/stfsy/go-api-kit/commit/d12025557b506431a11bdc7d7476f4243120a4c4))

## [1.11.0](https://github.com/stfsy/go-api-kit/compare/v1.10.0...v1.11.0) (2025-10-05)


### Features

* decode delete body too if content-length &gt; 0 ([2ccc70d](https://github.com/stfsy/go-api-kit/commit/2ccc70d6d9a535b7055c34a0bc2dfe60a2931ac0))

## [1.10.0](https://github.com/stfsy/go-api-kit/compare/v1.9.0...v1.10.0) (2025-09-24)


### Features

* error details have lowercase keys ([a64299d](https://github.com/stfsy/go-api-kit/commit/a64299dd08226d826acb459dbeb240634a495461))

## [1.9.0](https://github.com/stfsy/go-api-kit/compare/v1.8.0...v1.9.0) (2025-09-23)


### Features

* add functions to get safe header values ([26b9663](https://github.com/stfsy/go-api-kit/commit/26b96634bba48042660c24c9b1845ef00be61a19))
* add helper to create error details ([57855c4](https://github.com/stfsy/go-api-kit/commit/57855c4f9b853141658755f608c55958e3e413db))

## [1.8.0](https://github.com/stfsy/go-api-kit/compare/v1.7.0...v1.8.0) (2025-09-21)


### Features

* **server:** require cors config ([369aa0b](https://github.com/stfsy/go-api-kit/commit/369aa0b7e9ec1dcadd972040d491d1013fce8108))

## [1.7.0](https://github.com/stfsy/go-api-kit/compare/v1.6.0...v1.7.0) (2025-09-19)


### Features

* allow http 1.1 or higher ([7179f7e](https://github.com/stfsy/go-api-kit/commit/7179f7e05b53344bd39cea73ae2711c4a0b87801))
* also support chunked encodings ([b8d5e97](https://github.com/stfsy/go-api-kit/commit/b8d5e97ca84158e20019b569909a8219f1cd7195))
* disallow unknown json fields ([fcd6d5f](https://github.com/stfsy/go-api-kit/commit/fcd6d5ff955bc73ad8c0f2dc4ff9ff8b16d3f65d))
* ensure struct field lookups are cached ([dac68b1](https://github.com/stfsy/go-api-kit/commit/dac68b15a1c47069584907554692cfe41b269703))


### Bug Fixes

* fix potential memory leak ([6a6fd3c](https://github.com/stfsy/go-api-kit/commit/6a6fd3c3496da5403ca14e4144a97eb052f2de91))

## [1.6.0](https://github.com/stfsy/go-api-kit/compare/v1.5.0...v1.6.0) (2025-09-13)


### Features

* check headers of put and patch requests too ([61fc4e6](https://github.com/stfsy/go-api-kit/commit/61fc4e6c945c6d69a3bb418442816b64371c015b))
* update error handling, load server config ([3b2c1e6](https://github.com/stfsy/go-api-kit/commit/3b2c1e6b526aa37bc4f376e87485eddea48a324f))


### Performance Improvements

* use map to store security headers ([24d08a4](https://github.com/stfsy/go-api-kit/commit/24d08a481f724095874e1bbe3b164b6b2803657a))

## [1.5.0](https://github.com/stfsy/go-api-kit/compare/v1.4.0...v1.5.0) (2025-09-11)


### Features

* add http 1.1 only middleware ([d6569ba](https://github.com/stfsy/go-api-kit/commit/d6569bab6e77b6040193c3dfd972378261fda780))
* add middleware checking content length header ([fd60332](https://github.com/stfsy/go-api-kit/commit/fd60332fd043672e5f32700ce3c96258bb312b13))
* allow customizing the middleware stack ([ed776e1](https://github.com/stfsy/go-api-kit/commit/ed776e17f8bdc5412939ff935c04fd5f31f17cbd))

## [1.4.0](https://github.com/stfsy/go-api-kit/compare/v1.3.0...v1.4.0) (2025-08-24)


### Features

* add new method to marshall and respond with struct ([3b3003c](https://github.com/stfsy/go-api-kit/commit/3b3003ce8d3b1b48421e90d32626fba3c6db0722))

## [1.3.0](https://github.com/stfsy/go-api-kit/compare/v1.2.0...v1.3.0) (2025-08-24)


### Features

* add struct validator ([430ec5a](https://github.com/stfsy/go-api-kit/commit/430ec5a370b2f97a2d256a299fc260a5fbb5af53))
* application/problem+json content type for all errors ([e511006](https://github.com/stfsy/go-api-kit/commit/e5110062bcbd4e72efc7bdf13c4327164a250806))
* enable request body validation ([fd075f1](https://github.com/stfsy/go-api-kit/commit/fd075f1563dc30ca1fafb9911a7559db01afbb94))

## [1.2.0](https://github.com/stfsy/go-api-kit/compare/v1.1.0...v1.2.0) (2025-08-23)


### Features

* add validating handler ([782458a](https://github.com/stfsy/go-api-kit/commit/782458a2d3805c6ea906909046e2bf8820523192))

## [1.1.0](https://github.com/stfsy/go-api-kit/compare/v1.0.0...v1.1.0) (2025-08-21)


### Features

* add max body length middleware ([186c36a](https://github.com/stfsy/go-api-kit/commit/186c36a056670c22a87061b05dc4e49d4f427bd3))
* add simple liveness handler ([e257183](https://github.com/stfsy/go-api-kit/commit/e25718342242c5988caa351b5a20e764e73789b1))
* **server:** add connection timeouts ([cb3398c](https://github.com/stfsy/go-api-kit/commit/cb3398cd63170aa89164b65999f30efb97f51cef))

## 1.0.0 (2025-03-29)


### Features

* enable creation of multiple server instances with different ports ([823f17b](https://github.com/stfsy/go-api-kit/commit/823f17bd122b815d303940ecac54fd1c383619d0))
* **logger:** always log component ([976154f](https://github.com/stfsy/go-api-kit/commit/976154f2b670b0f46adb8b708cd6a3eeafa7de1b))
* **server:** on windows listen to localhost explicitly ([2138d6e](https://github.com/stfsy/go-api-kit/commit/2138d6e4d85364ff57b68f9ec903405f37c61716))
