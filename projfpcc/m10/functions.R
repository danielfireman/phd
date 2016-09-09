memData <- function(nCores, id, type) {
  p <- paste(gsub("nCores_", nCores, "../data/nCores_cores/server/memory_"), id, ".", type, sep = "")
  m <- merge(
    read.csv(paste(p, ".init.csv", sep = "")),
    read.csv(paste(p, ".max.csv", sep = "")),
    by = "t",
    suffixes = c(".init", ".max"))
  m <- merge(m, read.csv(paste(p,".committed.csv", sep = "")), by = "t")
  m <- merge(m, read.csv(paste(p, ".used.csv", sep = "")),
                by = "t",
                suffixes = c(".commited", ".used"))
  colnames(m) <- c("t", "init", "max", "commited", "used")
  return(m)
}

threadData <- function(nCores, id) {
  p <- paste(gsub("nCores_", nCores, "../data/nCores_cores/server/threads_"), id, sep = "")
  t <- merge(
    read.csv(paste(p, ".count.csv", sep = "")),
    read.csv(paste(p, ".blocked.count.csv", sep = "")),
    by = "t",
    suffixes = c(".all", ".blocked"))
  t <- merge(t, read.csv(paste(p, ".timed_waiting.count.csv", sep = "")), by = "t")
  t <- merge(
    t,
    read.csv(paste(p, ".waiting.count.csv", sep = "")),
    by = "t",
    suffixes = c(".timed_waiting", ".waiting"))
  colnames(t) <- c("t", "all", "blocked", "waiting", "timed_waiting")
  t <- merge(
    t,
    read.csv(paste(p, ".runnable.count.csv", sep = "")),
    by = "t",
    suffixes = c(".waiting", ".runnable"))
  colnames(t) <- c("t", "all", "blocked", "waiting", "timed_waiting", "runnable")
  t["prop.threads.timed_waiting"] <- t$timed_waiting / t$all
  t["prop.threads.waiting"] <- t$waiting / t$all
  t["prop.threads.blocked"] <- t$blocked / t$all
  t["prop.threads.runnable"] <- t$runnable / t$all
  return(t)
}

getCpuData <- function(nCores, id) {
  p <- paste(gsub("nCores_", nCores, "../data/nCores_cores/server/cpu_"), id, sep = "")
  t <- read.csv(paste(p, ".time.csv", sep = ""))
  l <- read.csv(paste(p, ".load.csv", sep = ""))
  # convert from nanosecs to millisecs
  t["value"] <- t$value / 1000
  # It is cumulative, we would like to work with an instantaneous version
  t$value <- c(NA, t[2:nrow(t), 2] - t[1:(nrow(t)-1), 2])
  t <- merge(t, l, by = "t")
  colnames(t) <- c("t", "cpu.time", "cpu.load")
  return(t)
}

gcTime <- function(nCores, id, type) {
  p <- paste(gsub("nCores_", nCores, "../data/nCores_cores/server/gc_"), id, ".", type, sep = "")
  t <- read.csv(paste(p, ".time.csv", sep = ""))
  # It is cumulative, we would like to work with an instantaneous version
  t$value <- c(NA, t[2:nrow(t), 2] - t[1:(nrow(t)-1), 2])
  return(t)
}

getGcCpuData <- function(nCores, id) {
  cpu <- getCpuData(nCores, id)
  # GC
  gc_ms <- gcTime(nCores, id, "PS-MarkSweep")
  gc_s <- gcTime(nCores, id, "PS-Scavenge")
  gc <- merge(gc_s, gc_ms, by = "t")
  colnames(gc) <- c("t", "gc.minor.time", "gc.full.time")
  gc["gc.time.sum"] <- gc$gc.minor.time + gc$gc.full.time
  
  # GC e CPU
  gc_cpu <- merge(gc, cpu)
  gc_cpu["prop.cpu.gc"] <- gc_cpu$gc.time.sum/gc_cpu$cpu.time
  return(gc_cpu)
}

plotMem <- function(df, yLab) {
  print(ggplot(size = 1) +
          geom_line(data = df, aes(x = t/ 1000, y = commited / 1000000, color="Commited")) +
          geom_line(data = df, aes(x = t/ 1000, y = used / 1000000, color="Used")) +
          geom_line(data = df, aes(x = t/ 1000, y = max / 1000000, color="Max")) +
          ylab(yLab) +
          xlab("Time (s)") +
          scale_color_manual(values = COLOR_BLIND_PALETTE) +
          theme(
            legend.title=element_blank(),
            legend.position = c(0.13, 0.9),
            text = element_text(size=8),
            axis.text = element_text(size = 8)))
}

plotThreads <- function(df) {
  print(ggplot(size = 1) +
          geom_line(data = threads, aes(x = t/ 1000, y = all, color = "all")) +
          geom_line(data = threads, aes(x = t/ 1000, y = timed_waiting, color = "timed waiting")) +
          geom_line(data = threads, aes(x = t/ 1000, y = waiting, color = "waiting")) +
          geom_line(data = threads, aes(x = t/ 1000, y = blocked, color = "blocked")) +
          geom_line(data = threads, aes(x = t/ 1000, y = runnable, color = "runnable")) +
          ylab("Count") +
          xlab("Time (s)") +
          scale_color_manual(values = COLOR_BLIND_PALETTE) +
          theme(
            legend.title=element_blank(),
            legend.position = c(0.8, 0.6),
            text = element_text(size=8),
            axis.text = element_text(size = 8)))
}

plotThroughput <- function(df) {
  print(
    ggplot(size = 1) +
      geom_line(data = df, aes(x = t/ 1000, y = mean_rate, color = "throughput")) +
      ylab("Throughput (Req/Sec)") +
      xlab("Time (s)") +
      scale_color_manual(values = COLOR_BLIND_PALETTE) +
      theme(
        legend.title=element_blank(),
        legend.position = c(0.1, 0.9),
        text = element_text(size=8),
        axis.text = element_text(size = 8))
  )
}


plotResponseTime <- function(df) {
  print(
    ggplot(size = 1) +
      geom_line(data = df, aes(x = t/ 1000, y = mean, color = "response time")) +
      ylab("Response Time (ms)") +
      xlab("Time (s)") +
      scale_color_manual(values = COLOR_BLIND_PALETTE) +
      theme(
        legend.title=element_blank(),
        legend.position = c(0.1, 0.9),
        text = element_text(size=8),
        axis.text = element_text(size = 8))
  )
}
