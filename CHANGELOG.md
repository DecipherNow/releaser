## 0.3.3 (May 20, 2021)

## Added

- If a tag includes `rc.-`, release only the fully qualified version name

## Changed

- Upgrade to Go 1.16
- Use Go modules

## 0.3.2 (Jan 8th, 2018)
### Added
- Push to normalized Docker registry names

## 0.3.1 (Jan 8th, 2018)
### Fixed
- Build binary without CGO to run on alpine

## 0.3.0 (Jan 8th, 2018)
### Added
- Option to add release assets to existing release ID

### Fixed
- Bug in parsing github token


## 0.2.1 (October 16th, 2018)
### Fixed
- Bug in splitting source docker tag off before re-tagging
- Bug in parsing `githubToken` rather than `token`

## 0.2.0 (October 10th, 2018)
### Added
- Docker image tagging now accepts a suffix to append to symver tags

### Changed
- Multiple release assets can now be specified

## 0.1.0 (October 10th, 2018)
### Changed
- Release asset is now optional
- Update README with simple usage instructions

## 0.0.0 (Beginning of Time)
### Added
- Create and push docker images
- Create github release and uppload 1 asset
