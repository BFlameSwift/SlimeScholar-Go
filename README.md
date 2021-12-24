# Slime Scholar backend



初始化包 ： `go mod tidy`,下载与生成依赖 

运行 go run ./main.go git








# 致谢

虽然获取数据的过程一波三折，还是要感谢微软提供了其全部的开源数据，也是本网站最主要的数据来源 本学术网站的主要数据来源：Microsoft-Academic-Graph

同样感谢Open-Alex提供了Microsoft-Academic-Graph的部分备用数据，虽然在最终并未使用，但是想必以后的同学们想要获取高质量数据的话，


OpenAlex将会是未来的第一选择。（MAG已宣布在2021年末停止服务。）

下面推荐几款个人在筛选中觉得还不错的数据（排名按照个人使用情况以及推荐度排序）：





[OpenAlex](https://openalex.org/data-dump) : 或许是MAG 的最好替代品。

[SemanticScholar](https://www.semanticscholar.org/product/api) ： 个人觉得非常牛的数据了，（不算mag的话）。唯一的缺点个人觉得就是没有机构以及领域的数据。会议的数据个人觉得还是相对无关紧要的。** 数据每月更新：** 虽然大多优秀的网站都是如此。总数据量大概在2.1亿左右 400G?（2021年11月统计），就连超级全的mag也只有2.6亿。在数据量上可以随便秒杀一大多常规数据了。但也因为以上痛点个人还是最终选择了 mag
	
此外此网站还有一些比较好用的API在链接处。至于数据源则是在[此](https://api.semanticscholar.org/corpus/)

OAG：from [aminer](https://www.aminer.cn/oag-2-1)

[scigraph](https://sn-scigraph.figshare.com/articles/dataset/Dataset_GRID_Organizations_for_SciGraph/7376537) ：看起来还算不错。应该涵盖了大部分的CSpaper

[Unpaywall](https://unpaywall.org/products/snapshot) :较优秀的开源数据网站，填写表格即可拿到数据。 解压后大概在130G。不过此时我已有mag便没有过多了解，但是很多数据网站均有提到此网站



最后来一个百家汇，是个列举开源数据网站的网站，虽然以上只有scigraph是我在这里发现的，但是不可否认这个网站还是比较全面的

[https://shubhanshu.com/awesome-scholarly-data-analysis/#networks](https://shubhanshu.com/awesome-scholarly-data-analysis/#networks)
