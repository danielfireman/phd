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
Lets say you would like to process one day worth of Github data in R.

* First pick your archive files (1 per hour) from [Github Archive](https://www.githubarchive.org/):

```{bash}
wget http://data.githubarchive.org/2015-01-01-{0..23}.json.gz
```

Then process everything and create a CSV (144090 lines):

```{bash}
$ output=2015-01-01.csv && \
rm $output && \
echo "EventID,Type,CreatedAt,Actor.ID,Repo.ID" >> $output && \
for f in `ls`; do \
echo "Processing $f" && \
cat $f | archive2csv  >> $output; \
done

Processing 2015-01-01-0.json.gz
7702 lines converted.
Processing 2015-01-01-10.json.gz
5370 lines converted.
Processing 2015-01-01-11.json.gz
7387 lines converted.
Processing 2015-01-01-12.json.gz
8273 lines converted.
...
```
