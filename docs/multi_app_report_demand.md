## 增加糖果数据查询（包括发行数据）👌
- 查询每日糖果发行总数量（SIE）       table sugars.sugar
- 查询每日0点流通量，糖果发行前（SIE） table sugars.currency
- 查询每日0点定存快照数量（SIE）      table sugars.defi_saving
- 查询每日0点抵押数量（SIE）         table sugars.defi_plegde
- 查询每日0点销毁总数量（SIE）        table sugars.shop_used_sie
- 查询每日糖果平均增长率（SIE）       table sugars.avg_growth_rate
- 查询双币质押的总数量（SIE）         table sugars.secret_chain_pledge

## 增加基金数据查询
### 👌 查询当天的基金购买数据
- 商户号
- 数量

### 👌 查询当天的基金发放数据
- 数量 proxy defi_fund api /platform_snapshot/:date

## 👌 六、增加挖矿数据查询
- 查询每日0点，双币质押的数量（GAS、SIE）（生效、待生效）     /sec_chain/i/pledge/snapshot/:date 获取当日的快照，再合计
- 查询每日gas发放数量                                   table settlements.gas_reward_volume 

## 👌 七、增加none数据查询
- 查询none当日发放金额(todo how to summarize)

## 👌 八、增加共识数据查询
- 查询当日公识发放数据（两次） job服务器搞一个静态文件服务，提供访问
- 查询公识总链接数量