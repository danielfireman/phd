# Abstract
*problem? why is it important? what is new being presented? conclusions?*

Getting accurate performance metric models of a web system is very important to solve the problems about performance analysis, evaluation and capacity planning, etc. Because of the complexity of a web system (e.g. a lot of software and hardware are integrated in the system), using analytical modeling method without integrating performance testing process is not enough to get accurate metric models.

To integrate performance testing process and analytical modeling method in a systematic and repeatable way, the authors have built a web performance modeling process under the methodology of learning from data. The process integrates the performance testing and modeling as a whole and divides a performance modeling activity into several phases. To demonstrate the process, a performance scalability problem of a real web community system (www.igroot.com) is studied.

Experiments showed that result models can be used to demonstrate the performance requirement indexes, estimate the saturation and buckle points, estimate the performance tinning space and facilitate finding performance bottlenecks.

# My opinion

The paper is strongly based on *Web Application Scalability: A Model-Based Approach (Williams and Smith)*. It adds the idea of learning from data to build a process which integrates testing and analytical modeling, which is very interesting. Even though results look promising, they only tested in one web site.

Another point is that some inputs need to be determined by benchmarks (i.e. benchmark number of concurrent users.). How to select and configure those benchmarks still not covered by the article.

As the result of the paper is a process proposal. As it is intended to be generally applicable and only one example was shown, it was for me hard to understand how all the options to profile, infer, validate and so on were selected (practical aspects).

# How the article relates to my work

My work is about performance modeling, in particular in the context of cloud computing (this bit, not fully related to the article). Another two important relations: i) I wasn't considering the idea of a proposing a process, which could be good and ii) the idea of learning from data seemed promising and I definitely should look more into that.