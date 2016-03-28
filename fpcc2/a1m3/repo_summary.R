source("data_load.R")


summary(repoCount.2012$n)
summary(repoCount.2013$n)
summary(repoCount.2014$n)
summary(repoCount.2015$n)

head(repoCount.2012, n=5)
head(repoCount.2013, n=5)
head(repoCount.2014, n=5)
head(repoCount.2015, n=5)

skewness(repoCount.2012$n)
skewness(repoCount.2013$n)
skewness(repoCount.2014$n)
skewness(repoCount.2015$n)

mode(repoCount.2012$n)
mode(repoCount.2013$n)
mode(repoCount.2014$n)
mode(repoCount.2015$n)