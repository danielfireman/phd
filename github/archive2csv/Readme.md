# Convert Github Archive to CSV files

This is a small tool that aims to help analysis of [Github Archive](https://www.githubarchive.org/) files by converting the JSON records to CSV. In my case, I am going to use [R](https://www.r-project.org/) to explore the result data set.

Please feel free to open issues. Pull requests are welcomed!

## Installation

```{bash}
# If you don't have Golang installed you could use [GVM](https://github.com/moovweb/gvm) to make the installation easier
$ go get github.com/gocarina/gocsv
$ go get github.com/danielfireman/phd/github/archive2csv
```

If you don't want to install Go, feel free to send me gentle e-mail saying your platform and I will make the binary avaiable.

## Usage

First you pick your archive file at [Github Archive](https://www.githubarchive.org/), lets say ~/githubarchive/2015-03-01-0.json.gz. Then run:

```{bash}
$ archive2csv --archive_path=~/githubarchive/2015-03-01-0.json.gz
Processing ~/githubarchive/2015-03-01-0.json.gz. Output will be written to ~/githubarchive/2015-03-01-0.csv
Conversion completed.
```
