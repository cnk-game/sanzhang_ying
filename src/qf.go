package main
import (
	"regexp"
	"fmt"
	"encoding/json"
	"game/domain/user"
	"flag"
	"strings"
	"net/http"
	"io"
	"io/ioutil"
	"strconv"
)

var s = `I0605 14:57:35.704541 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1749779","order":"df577327-29ae-4510-d","price":10,"payType":1,"payCode":100013,"state":"success","time":"2015-6-5 14:55:53","gameOrder":"55699c301d4bd4280907e019:6a22207f-8f02-4784-8e25-70b441b168a5","sign":"62fbd75b20cd5e14432976362eb43df4"}
I0605 15:20:16.719638 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1453298","order":"1a68c051-a33c-80e8-a","price":10,"payType":1,"payCode":100013,"state":"success","time":"2015-6-5 15:19:50","gameOrder":"5571203b1d4bd4280908ed65:96d6aeb8-5696-408b-ae84-09475451e10c","sign":"a826d441e3b71d0747a0fefe19ab8d5c"}
I0605 15:22:47.187175 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1453298","order":"1e36bd77-e874-baf6-f","price":10,"payType":1,"payCode":100013,"state":"success","time":"2015-6-5 15:22:25","gameOrder":"5571203b1d4bd4280908ed65:69dd3157-db01-4e45-b6be-bc27116894b2","sign":"2e57381c3d901460b13a6bce191fd5df"}
I0605 15:28:29.967349 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1453298","order":"3326408c-64de-e9aa-f","price":6,"payType":1,"payCode":100012,"state":"success","time":"2015-6-5 15:27:52","gameOrder":"5571203b1d4bd4280908ed65:9316400e-d6a7-47fa-8ee6-02ae8e482368","sign":"7a3b9331c75e34d4ca4fac9fd866074a"}
I0605 15:35:14.065925 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1453298","order":"86a67329-ed42-055b-5","price":10,"payType":1,"payCode":100013,"state":"success","time":"2015-6-5 15:34:18","gameOrder":"5571203b1d4bd4280908ed65:f9189167-2e2c-4670-a8d6-adcb675fb975","sign":"24ef0d3934bc44d62a60ff779fd5090f"}
I0605 16:45:59.705937 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1830598","order":"0aa70ab8-c17a-803d-9","price":10,"payType":1,"payCode":100013,"state":"success","time":"2015-6-5 16:45:39","gameOrder":"5570f84a1d4bd4280908e750:30617eb7-917b-4d98-a89a-ac38c51e90b9","sign":"d8f327c33a53f6becd2881d696e8042c"}
I0605 16:59:15.271317 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1327658","order":"9e67bd19-e86c-b55b-5","price":10,"payType":7,"payCode":100013,"state":"success","time":"2015-6-5 16:58:51","gameOrder":"556fac771d4bd4280908b5e1:6974e9c0-c31b-45a5-b79d-e199c6bc222a","sign":"a8b9be834db79f0a485615576a724aa5"}
I0605 17:52:27.283273 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1749779","order":"5f4e7603-e9cb-2022-8","price":10,"payType":1,"payCode":100013,"state":"success","time":"2015-6-5 17:50:33","gameOrder":"55699c301d4bd4280907e019:e4bb3d80-310c-403b-9559-b9df8912a1ce","sign":"40d306a2ef042bcd83c1e72c03183b60"}
I0605 18:16:42.018875 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"874808","order":"3bcb8bb7-1e54-3e57-8","price":10,"payType":1,"payCode":100013,"state":"success","time":"2015-6-5 18:16:23","gameOrder":"555f86081d4bd4280905c5d1:0ffce7e3-00bc-4666-b3a3-b8c789f0abd4","sign":"35e80967edb572dce19545a20b387eca"}
I0605 18:19:53.418195 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1934501","order":"085e50ba-cb20-7b30-9","price":10,"payType":1,"payCode":100013,"state":"success","time":"2015-6-5 18:18:13","gameOrder":"55706a4e1d4bd4280908dc24:0c02cf76-dd8d-44f6-af52-7d6c6bc569ff","sign":"94040757c0801dc6909a489d27a84b0b"}
I0605 18:21:29.521650 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1934501","order":"c1a49590-0a30-2687-1","price":6,"payType":1,"payCode":100012,"state":"success","time":"2015-6-5 18:19:52","gameOrder":"55706a4e1d4bd4280908dc24:6d2067b4-a921-418a-b687-a71a89899e29","sign":"d6800839b078ff0048ceaed90ec817c6"}
I0605 18:34:47.396701 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1327658","order":"f008f8f6-4dda-24e2-6","price":10,"payType":3,"payCode":100013,"state":"success","time":"2015-6-5 18:34:0","gameOrder":"556fac771d4bd4280908b5e1:875d1bb7-014e-42a8-bc5e-985e6642475e","sign":"2f5e23f3b66a1248a3092a376f9b2f48"}
I0605 18:42:29.303049 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1624665","order":"6f248c38-365d-43f4-5","price":6,"payType":2,"payCode":100012,"state":"success","time":"2015-6-5 18:42:13","gameOrder":"554ab6ef1d4bd4280901e480:2cab6056-376f-4388-b5ce-435b30f029b7","sign":"72e8c976603f375384ea89a9b8ae956d"}
I0605 18:57:56.824281 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1885993","order":"ec8fbf01-d993-991e-4","price":10,"payType":1,"payCode":100013,"state":"success","time":"2015-6-5 18:50:18","gameOrder":"557036191d4bd4280908cfa7:a5f5b2e0-0610-4628-bf1b-206ddc091ab1","sign":"e47b70c608fcae3d742a83f22e86181c"}
I0605 20:38:27.460239 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"927239","order":"c98ba936-40e5-f71b-3","price":6,"payType":1,"payCode":100012,"state":"success","time":"2015-6-5 20:38:3","gameOrder":"55715cd01d4bd4280908f897:871a36fa-d118-47c8-becb-82905b1e904d","sign":"7c1d696b3978394fa765514c8d7426d5"}
I0605 20:52:19.566435 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1684724","order":"b83ab9fb-cea7-3c01-4","price":6,"payType":1,"payCode":100012,"state":"success","time":"2015-6-5 20:52:1","gameOrder":"5552c1be1d4bd42809034227:412eaddd-4f06-4113-998a-1ff66e89c216","sign":"f3025f97f57fdefb65484907006434c6"}
I0605 21:04:42.744857 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1935315","order":"da523c8a-3da2-ca81-d","price":10,"payType":2,"payCode":100013,"state":"success","time":"2015-6-5 21:3:50","gameOrder":"557183c11d4bd42809090072:274ebb8b-5ed1-49d5-993e-4afba0d915bd","sign":"6f0be12548e8756e612aabc16ac8832a"}
I0605 21:57:39.809253 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1684724","order":"856d4982-8614-11c5-e","price":6,"payType":1,"payCode":100012,"state":"success","time":"2015-6-5 21:57:24","gameOrder":"5552c1be1d4bd42809034227:8d8e8cd7-55ce-4745-82a8-0c8b71b28ede","sign":"48a282af6bd61bb0407e1bf25df0a2a1"}
I0605 21:58:39.225877 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1684724","order":"9667c872-b3c0-eff3-8","price":6,"payType":1,"payCode":100012,"state":"success","time":"2015-6-5 21:58:26","gameOrder":"5552c1be1d4bd42809034227:548fd6fd-1fe5-4461-b6e6-4a57a8dae0b2","sign":"45768bcd835f4d68eae62230e7377492"}
I0605 22:14:55.975221 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1604328","order":"2f26f3ef-86cc-a706-0","price":10,"payType":1,"payCode":100013,"state":"success","time":"2015-6-5 22:14:15","gameOrder":"5548793c1d4bd428090183b8:b6d70a15-ef7c-457e-93c3-ee66839c0c70","sign":"33f429f1f5fd514b48b9de74a4f1e883"}
I0606 00:31:58.644567 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"170254","order":"62b205f7-97ba-a6fe-9","price":10,"payType":3,"payCode":100013,"state":"success","time":"2015-6-6 0:31:39","gameOrder":"551c16f11d4bd434e500687d:c72c7994-10f5-4917-8d8d-49b029bc7bf4","sign":"f7460212cfefb950aa0a15ef12f4e710"}
I0606 03:24:38.133216 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1604328","order":"e1226a20-744c-9752-5","price":6,"payType":1,"payCode":100012,"state":"success","time":"2015-6-6 3:24:25","gameOrder":"5548793c1d4bd428090183b8:16cee112-4248-4bf2-ba20-36fdd757cf1d","sign":"62ab7bacc1d8acbb88c4b8bd7d6aa7fe"}
I0606 03:30:22.175981 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1604328","order":"0b70be86-2b85-e314-0","price":6,"payType":1,"payCode":100012,"state":"success","time":"2015-6-6 3:25:2","gameOrder":"5548793c1d4bd428090183b8:e372a1e0-aff0-4cdd-9f25-26e622d77a81","sign":"19a2c4e29386cb36d4e64786197d41b3"}
I0606 03:50:04.200135 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1624665","order":"b2516444-a7cd-7ce0-2","price":10,"payType":2,"payCode":100013,"state":"success","time":"2015-6-6 3:44:41","gameOrder":"554ab6ef1d4bd4280901e480:07bf42ff-7575-43bc-b30e-02813b446749","sign":"5d0017b8caacc0d30c9edd739330e0c3"}
I0606 11:53:32.788939 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"460068","order":"dcde489e-7d50-6db9-e","price":6,"payType":1,"payCode":100012,"state":"success","time":"2015-6-6 11:53:15","gameOrder":"556883371d4bd4280907bc61:9198ecb7-a374-42ba-be03-03dca3a94617","sign":"00a40521954292780afd6ac8d8724db8"}
I0606 11:54:19.325553 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"460068","order":"61144af1-4141-d13d-d","price":6,"payType":1,"payCode":100012,"state":"success","time":"2015-6-6 11:53:59","gameOrder":"556883371d4bd4280907bc61:a7ab614a-21c3-406d-87d3-ef12cf4068ff","sign":"fcc0adb6a2a0de9928f53a6527167804"}
I0606 11:54:51.101913 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"460068","order":"b56c69fa-e18e-cb48-6","price":6,"payType":1,"payCode":100012,"state":"success","time":"2015-6-6 11:54:31","gameOrder":"556883371d4bd4280907bc61:35bb8767-dbb1-41e2-a8f7-71efe887a7bb","sign":"18edb6900367539f09f8e594582b0e16"}
I0606 12:06:24.469874 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1418733","order":"cb43483e-a22f-ed46-8","price":10,"payType":1,"payCode":100013,"state":"success","time":"2015-6-6 12:6:8","gameOrder":"557021351d4bd4280908ca22:71f6d3cf-0d58-465a-add4-9b27f8f076ed","sign":"0ae444915ed0f7ed509b0ad7e1ac3db3"}
I0606 12:17:34.202615 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1418733","order":"da91a8cc-ed04-ee85-e","price":10,"payType":1,"payCode":100013,"state":"success","time":"2015-6-6 12:17:20","gameOrder":"557021351d4bd4280908ca22:64e7b908-125f-4af3-a0c5-edbf23fb5ecd","sign":"6f59b3e633ae779718aa34798d3d2436"}
I0606 12:55:17.667401 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1847677","order":"08ab5170-8a29-72ee-6","price":10,"payType":3,"payCode":100013,"state":"success","time":"2015-6-6 12:54:57","gameOrder":"556a7b091d4bd4280907fac7:9dab6a96-cebb-4c55-a16b-d3c13a338030","sign":"10d94d10226dd28a4be50c43d088de47"}
I0606 14:48:11.853831 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1157853","order":"c5ab91b6-a2c9-870a-a","price":6,"payType":3,"payCode":100012,"state":"success","time":"2015-6-6 14:47:43","gameOrder":"5562f6581d4bd4280906b882:c00a2ee6-c41d-48fd-82d7-820c85ff33d5","sign":"72c79c991b78060cc164e9296f3ab433"}
I0606 16:02:36.943981 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1308123","order":"c8396f75-5eca-e7a2-0","price":10,"payType":7,"payCode":100013,"state":"success","time":"2015-6-6 16:2:11","gameOrder":"556ffc381d4bd4280908c407:3a20ea49-310b-443b-b50a-a312e8d1c3ba","sign":"ef2515cb5afed54a3613ada6442a973b"}
I0606 18:23:12.199354 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1027040","order":"bfdbbe3a-a9f1-4fc3-7","price":10,"payType":1,"payCode":100013,"state":"success","time":"2015-6-6 18:22:52","gameOrder":"55380d271d4bd4659c017159:1bf31384-6d90-43d1-838b-22c2f5a133a7","sign":"97d0cef95118c20c1b33ab83011634ad"}
I0606 18:29:28.237431 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1098425","order":"d088a521-196c-a046-a","price":10,"payType":1,"payCode":100013,"state":"success","time":"2015-6-6 18:29:12","gameOrder":"555353101d4bd428090360db:a1e612ee-7c4f-45d5-b8a2-63ef890c7b08","sign":"e59d75a9fc0dae5964a920b55a8260c8"}
I0606 18:29:52.973422 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1098425","order":"db536c4d-19f4-466f-b","price":10,"payType":1,"payCode":100013,"state":"success","time":"2015-6-6 18:29:38","gameOrder":"555353101d4bd428090360db:34ac8686-6a90-4e77-a44c-30fd8d6decf0","sign":"770ef442e278b3e5175b836438a8f36b"}
I0606 18:51:07.171948 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"170254","order":"206e6e93-500e-ab13-0","price":10,"payType":3,"payCode":100013,"state":"success","time":"2015-6-6 18:50:48","gameOrder":"551c16f11d4bd434e500687d:0a568f21-72b1-41ad-9018-ad0ad2d55f02","sign":"cca9930ab0dfe2090c233c61bf3b750b"}
I0606 19:18:28.213376 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"566780","order":"59705b10-2966-cba8-2","price":6,"payType":1,"payCode":100012,"state":"success","time":"2015-6-6 19:18:20","gameOrder":"5572d6901d4bd428090933cd:ebb4e475-5830-4796-97d4-a87c167f9907","sign":"63ae0a4131bf59ebff34b64b84fa1390"}
I0606 19:20:09.852848 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1308123","order":"2efcc66c-7633-3332-7","price":10,"payType":7,"payCode":100013,"state":"success","time":"2015-6-6 19:19:58","gameOrder":"556ffc381d4bd4280908c407:0df7a2c1-2661-45e0-aabb-cdd9b8f62dab","sign":"eb098766d9bf5af83349456b575d64e6"}
I0606 20:53:41.185638 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1830598","order":"25dc7470-e725-26a4-1","price":6,"payType":1,"payCode":100012,"state":"success","time":"2015-6-6 20:53:28","gameOrder":"5570f84a1d4bd4280908e750:5d370d7f-d14d-42ae-8ffe-333b01873de6","sign":"85e746302244a819d2a41bec19b16764"}
I0606 20:55:17.978713 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1904945","order":"5a530550-765f-3030-3","price":6,"payType":3,"payCode":100012,"state":"success","time":"2015-6-6 20:55:7","gameOrder":"556dba4c1d4bd428090879da:dcfa31fe-8bb5-4eab-a147-8ea706193329","sign":"ef27ea67cb57ee4ad578acd3e4f7c5f1"}
I0606 20:55:26.897605 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1904945","order":"37745449-4bae-5b5f-e","price":10,"payType":3,"payCode":100013,"state":"success","time":"2015-6-6 20:55:20","gameOrder":"556dba4c1d4bd428090879da:87d2129a-79c4-4909-9a56-48ee187b00ff","sign":"2e3fee66f45d77ee5fe06365da846fc9"}
I0606 21:00:56.555368 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1830598","order":"808bea08-4f0e-016e-3","price":10,"payType":1,"payCode":100013,"state":"success","time":"2015-6-6 21:0:45","gameOrder":"5570f84a1d4bd4280908e750:ed36f9f1-09ac-4685-9e17-fd070b026cc8","sign":"f0581ff24cbe10218a83303c0598d113"}
I0606 21:02:21.156999 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1904945","order":"51cfc111-3b70-9bea-d","price":10,"payType":3,"payCode":100013,"state":"success","time":"2015-6-6 21:1:52","gameOrder":"556dba4c1d4bd428090879da:34e883cc-6e26-44d9-9431-53549040feaf","sign":"1c18e36b11d3baf5b55f400344b3396d"}
I0606 22:00:04.922897 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"709943","order":"3433e7b4-652a-138c-5","price":6,"payType":1,"payCode":100012,"state":"success","time":"2015-6-6 21:59:50","gameOrder":"55722f4b1d4bd428090915c0:c8dc5bb6-4f84-4ba0-a6a3-4496ade435c8","sign":"175ef3a3ccde9271080fd78b1a598e7b"}
I0606 22:05:47.610623 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"709943","order":"c99a3ac9-19b0-1f2b-8","price":10,"payType":1,"payCode":100013,"state":"success","time":"2015-6-6 22:5:36","gameOrder":"55722f4b1d4bd428090915c0:3e16da5f-5656-4837-95d5-0846c240749c","sign":"6c6fb4bba215be27cc3b046825263e45"}
I0606 23:14:02.353389 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1684724","order":"496c8c62-888b-08d5-1","price":6,"payType":1,"payCode":100012,"state":"success","time":"2015-6-6 23:13:45","gameOrder":"5552c1be1d4bd42809034227:9a0964e9-1251-4c1e-b0dd-f3416f92cbfe","sign":"e30ba806002e5fb08b0f21d4671064b1"}
I0606 23:14:43.740473 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1684724","order":"fcce88fa-50c6-a6ef-3","price":6,"payType":1,"payCode":100012,"state":"success","time":"2015-6-6 23:14:22","gameOrder":"5552c1be1d4bd42809034227:6df2b327-45e1-4142-b2a0-9552a839fa7f","sign":"a7ee0c654fad6b5f2f0b23030285a04f"}
I0606 23:16:17.398141 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1684724","order":"a645aff2-c49e-a569-1","price":6,"payType":1,"payCode":100012,"state":"success","time":"2015-6-6 23:15:57","gameOrder":"5552c1be1d4bd42809034227:04cec5a2-f2e4-4590-b619-24c8374ecff7","sign":"ab7092ad1b90b5bab21c418bfed208e2"}
I0606 23:29:01.848515 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"448915","order":"f0897275-45a7-ad39-e","price":6,"payType":7,"payCode":100012,"state":"success","time":"2015-6-6 23:28:48","gameOrder":"554a0a7a1d4bd4280901cda6:0aa12ebd-a54f-4f1d-9ef6-2c16f63802a5","sign":"a0330c831f47b5fc4728bdccbdf3a883"}
I0606 23:29:37.147307 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"448915","order":"b299cb42-edac-8867-3","price":10,"payType":7,"payCode":100013,"state":"success","time":"2015-6-6 23:29:26","gameOrder":"554a0a7a1d4bd4280901cda6:f64bf21d-ab51-419e-801c-e665ff5aeb46","sign":"91468308aef30eb7ea12a6dbd982b168"}
I0607 00:02:30.997448 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1856827","order":"d3bd0dd5-eb25-34d1-4","price":6,"payType":1,"payCode":100012,"state":"success","time":"2015-6-7 0:2:4","gameOrder":"556c6d2c1d4bd42809084ab7:81d3e705-6594-4fbd-8aae-3e9e4d69af98","sign":"b5f8b0a6782f69b06247f7690e20c5e5"}
I0607 00:35:15.897680 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1367322","order":"85cda717-a396-26e9-1","price":10,"payType":1,"payCode":100013,"state":"success","time":"2015-6-7 0:34:53","gameOrder":"552fe4361d4bd4659c00558b:4a23bfdb-ab33-46c5-b6e8-7711cc481da6","sign":"62ee9059ced2ba8c24734f22f74f2cf0"}
I0607 00:36:04.366519 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1367322","order":"2a369517-8e2c-988f-9","price":10,"payType":1,"payCode":100013,"state":"success","time":"2015-6-7 0:35:54","gameOrder":"552fe4361d4bd4659c00558b:17f71739-6228-45bc-b493-cbbe41f36386","sign":"fa00c17e404831c9166c01b58760fdb1"}
I0607 00:37:54.403364 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1367322","order":"4edc4174-63e2-70e8-6","price":10,"payType":1,"payCode":100013,"state":"success","time":"2015-6-7 0:37:10","gameOrder":"552fe4361d4bd4659c00558b:d1c79cbf-4a03-49f6-a7ae-2b22dbdf98af","sign":"6f8ceb17ff0f73e4e22590ffd835e8ec"}
I0607 00:47:30.180740 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1367322","order":"cf949bb0-5836-9f6e-a","price":6,"payType":1,"payCode":100012,"state":"success","time":"2015-6-7 0:42:0","gameOrder":"552fe4361d4bd4659c00558b:b76f111c-654a-4478-857a-dd97ab2e4e43","sign":"cc320b293d58beb343a35128d589879c"}
I0607 06:18:48.473343 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1914538","order":"cf6c7f19-b818-12b6-2","price":10,"payType":2,"payCode":100013,"state":"success","time":"2015-6-7 6:17:48","gameOrder":"556ed09d1d4bd428090899f7:eb9395c6-4652-448e-a3b5-bf11dfcebdac","sign":"4593883c80835bf9228ddac10346d1e9"}
I0607 06:21:08.598345 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1914538","order":"b28b3954-b9b5-3895-e","price":10,"payType":2,"payCode":100013,"state":"success","time":"2015-6-7 6:20:11","gameOrder":"556ed09d1d4bd428090899f7:2dee86bf-5792-4ad9-a3ec-6005c0ab7056","sign":"136a3cbd2e376071fab8dd43ad7282d2"}
I0607 06:26:29.233879 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1914538","order":"1db14c36-4987-3a6c-e","price":6,"payType":2,"payCode":100012,"state":"success","time":"2015-6-7 6:25:44","gameOrder":"556ed09d1d4bd428090899f7:ab2fd0df-16bc-49e4-a0e0-5e77485f1851","sign":"9e58699a1dd80362869d47edb36c50bc"}
I0607 06:40:56.096657 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1914538","order":"03e7c550-228e-9efc-f","price":10,"payType":2,"payCode":100013,"state":"success","time":"2015-6-7 6:39:42","gameOrder":"556ed09d1d4bd428090899f7:1a7d6500-2c76-410c-9e5a-bb9e1bae6bfe","sign":"3189f98a5bb815d83c4b8f9e7b9a971c"}
I0607 06:44:06.548008 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"328807","order":"796866e7-bad9-2cda-0","price":10,"payType":1,"payCode":100013,"state":"success","time":"2015-6-7 6:43:32","gameOrder":"551d16221d4bd434e5008ac7:74615ac5-e779-4018-8810-7823b839257b","sign":"12483f74383b6d0ee3d8a24863e7a7b5"}
I0607 06:44:37.150108 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"328807","order":"ff6fffd4-fb2a-4668-d","price":10,"payType":1,"payCode":100013,"state":"success","time":"2015-6-7 6:44:20","gameOrder":"551d16221d4bd434e5008ac7:b40d5ef3-294f-4930-a172-96a82c736fd4","sign":"8dac32f55d1d273afa302b125082e2d6"}
I0607 08:14:38.416335 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1892292","order":"d0046863-eaae-9136-6","price":10,"payType":1,"payCode":100013,"state":"success","time":"2015-6-7 8:14:10","gameOrder":"556b16041d4bd4280908190e:768b39ea-7bca-4da7-8093-ab863ff6247d","sign":"664de94a76a98ba029543b30e82331d3"}
I0607 09:22:39.384085 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1722200","order":"54db693d-5b64-4c17-1","price":6,"payType":1,"payCode":100012,"state":"success","time":"2015-6-7 9:22:9","gameOrder":"555729461d4bd4280904098c:f3ca7253-5152-4131-a31c-15e44dfc850d","sign":"23e82f1b7052461f633f8e7ad3b90a90"}
I0607 09:36:02.075321 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1956272","order":"65bdaf95-d63b-242c-0","price":10,"payType":1,"payCode":100013,"state":"success","time":"2015-6-7 9:33:52","gameOrder":"557397e31d4bd42809094b98:478a38f1-b6a6-4472-b0fa-e3975f1e5c0e","sign":"30e07cd2d187fd4a2fda26aa12da3bef"}
I0607 09:52:26.983982 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1689694","order":"88ebfc97-5bcb-5710-7","price":10,"payType":1,"payCode":100013,"state":"success","time":"2015-6-7 9:52:0","gameOrder":"555689b71d4bd4280903e62f:d216de4b-393e-4d42-83ad-a17529727a28","sign":"38526f4b3f35480a2dcf3672b3e7a728"}
I0607 10:29:49.513928 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1579696","order":"9e0e5fc2-a95c-8644-4","price":10,"payType":1,"payCode":100013,"state":"success","time":"2015-6-7 10:29:17","gameOrder":"554781a41d4bd42809016259:6b3dae6c-85f8-4c8e-85df-2ee91a84ebbe","sign":"1a1c462e6bc7dbc5205cac3e6dc23797"}
I0607 12:13:11.009800 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1327658","order":"bf4ebc27-77e4-6852-c","price":10,"payType":7,"payCode":100013,"state":"success","time":"2015-6-7 12:12:44","gameOrder":"556fac771d4bd4280908b5e1:53502a42-fec7-4638-8951-b581a9397730","sign":"ce220c8c44f5ff17a544c35055ddb4b5"}
I0607 14:13:25.405121 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1803875","order":"3b25ad4a-7826-00b8-0","price":10,"payType":2,"payCode":100013,"state":"success","time":"2015-6-7 14:12:50","gameOrder":"556060d31d4bd4280906092d:b60c5e76-768c-434c-93a7-37fa8a787e48","sign":"f248eb37b187ccc2a3e67412eea6eb17"}
I0607 14:24:54.230555 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1031227","order":"98c57dac-66a1-495a-8","price":10,"payType":7,"payCode":100013,"state":"success","time":"2015-6-7 14:22:43","gameOrder":"5573084a1d4bd42809093f6c:c2ddb8d2-6625-4169-8f39-63527356b441","sign":"957fc441a177fd61bcd3cab92560f5ab"}
I0607 14:27:03.519408 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1031227","order":"a162c634-e0dc-31b8-4","price":10,"payType":7,"payCode":100013,"state":"success","time":"2015-6-7 14:26:40","gameOrder":"5573084a1d4bd42809093f6c:85932a8b-8eba-4e9b-af45-a181c7f66257","sign":"9376a84e7adb7de41a47b88d346025d4"}
I0607 14:57:16.498620 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1763163","order":"13feeeb7-c77a-f3f8-a","price":10,"payType":2,"payCode":100013,"state":"success","time":"2015-6-7 14:56:53","gameOrder":"5573eacd1d4bd42809095b73:769e5009-8e29-4bb5-a009-09baae8c9136","sign":"a65f9c0861627c85d46543a49219921c"}
I0607 15:12:28.314021 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1324033","order":"62d44e58-87d7-c753-f","price":10,"payType":1,"payCode":100013,"state":"success","time":"2015-6-7 15:12:0","gameOrder":"5573e5831d4bd42809095aa7:b42eb16b-4291-4ef9-a9ab-96c79bc479de","sign":"eb575d09b53c893a10adbf5731fea95a"}
I0607 15:22:28.482116 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1031227","order":"d9c794a4-910b-e753-7","price":10,"payType":7,"payCode":100013,"state":"success","time":"2015-6-7 15:21:59","gameOrder":"5573084a1d4bd42809093f6c:feb5124b-87f7-4ad1-a885-d061a52ebfb3","sign":"d5ad88be2b0c1dca73320d201e702cd5"}
I0607 15:29:57.588486 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1885993","order":"59182eab-a5e3-98bd-3","price":10,"payType":1,"payCode":100013,"state":"success","time":"2015-6-7 15:22:24","gameOrder":"557036191d4bd4280908cfa7:65bfabd0-9419-44cd-8219-893af4e958ab","sign":"50384e77c1104424ef1b0d6c4872f0f0"}
I0607 15:53:22.627983 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1964201","order":"a8681afc-6477-5cb6-6","price":6,"payType":2,"payCode":100012,"state":"success","time":"2015-6-7 15:52:55","gameOrder":"5573f5e21d4bd42809095d19:7204c13c-720d-4024-9042-266204f74cd9","sign":"6a9a0fee47fbd58edc0a4afe2c13b71b"}
I0607 15:53:42.282655 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1964201","order":"b3338911-706f-84fc-8","price":6,"payType":2,"payCode":100012,"state":"success","time":"2015-6-7 15:53:24","gameOrder":"5573f5e21d4bd42809095d19:d392c3da-8420-4cfe-9f9d-5a95f2c04461","sign":"f268ffcddffcd5aacdf355b420019c0a"}
I0607 15:56:59.177362 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1844188","order":"88c2ead0-9208-f549-e","price":10,"payType":1,"payCode":100013,"state":"success","time":"2015-6-7 15:56:17","gameOrder":"5566e1b01d4bd42809077997:40f2f865-8b4e-44ac-8eae-70d052f123a2","sign":"59303941e5a1a6d4729a0eff9905a159"}
I0607 16:07:32.189261 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1844188","order":"cee9ce80-2c85-b397-6","price":10,"payType":1,"payCode":100013,"state":"success","time":"2015-6-7 16:6:57","gameOrder":"5566e1b01d4bd42809077997:d115378f-14a2-47e9-b309-c3c5c7d99775","sign":"4289d58671bc26b3c605a50a1d67fd98"}
I0607 16:17:14.453314 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1548744","order":"1b491abd-888e-2a59-6","price":6,"payType":2,"payCode":100012,"state":"success","time":"2015-6-7 16:16:46","gameOrder":"556d962b1d4bd42809087098:4ee42390-6d53-47c3-aa09-3dd3c80e1acf","sign":"d0a58e723b41067a7902c33af007b93c"}
I0607 16:48:48.919550 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1951716","order":"2b6ad5c0-0cf0-7d00-2","price":10,"payType":1,"payCode":100013,"state":"success","time":"2015-6-7 16:48:19","gameOrder":"557404cc1d4bd42809095f78:9202630e-d478-4646-9e89-a938f2f212e9","sign":"bcc687e436180174534207f4ebcc9cb1"}
I0607 16:52:26.158080 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1951716","order":"196a64a8-5f0d-574e-3","price":10,"payType":1,"payCode":100013,"state":"success","time":"2015-6-7 16:52:8","gameOrder":"557404cc1d4bd42809095f78:8aab13b9-b01a-47c3-9662-080dd9707f65","sign":"9e76d0b77483a08918280757501a0135"}
I0607 17:19:31.072369 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1951716","order":"5c5137a8-f243-7a60-3","price":6,"payType":1,"payCode":100012,"state":"success","time":"2015-6-7 17:19:11","gameOrder":"557404cc1d4bd42809095f78:a3fde994-2906-4578-b84f-cb554eccf818","sign":"4d7b8df4216b1cfb0f3ebefda45ba104"}
I0607 17:22:02.159607 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1951716","order":"f8857d67-8b92-1b4d-d","price":6,"payType":1,"payCode":100012,"state":"success","time":"2015-6-7 17:21:38","gameOrder":"557404cc1d4bd42809095f78:218f7e20-6500-49c2-869b-834d61b38ce4","sign":"4e47b7953e40fe0c4f8e79c1db5b0904"}
I0607 18:48:47.671241 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1467152","order":"d9f0d3d4-a41c-89fe-5","price":6,"payType":1,"payCode":100012,"state":"success","time":"2015-6-7 18:48:21","gameOrder":"555093a31d4bd4280902f259:41496eaf-07fa-471b-853e-67fd11ca64ed","sign":"28f3cf2a8d49aadfeeefcdef22225721"}
I0607 21:06:04.907200 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1466572","order":"4457e7ed-95a1-d24c-8","price":10,"payType":1,"payCode":100013,"state":"success","time":"2015-6-7 21:5:23","gameOrder":"5566a4741d4bd42809077099:4516d2af-53f4-46ad-8e56-817b1c4f3e6b","sign":"ea96c91e02448908996085979d2b8f40"}
I0607 21:23:21.356303 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"56266","order":"db7d5b36-5121-2f2e-8","price":10,"payType":1,"payCode":100013,"state":"success","time":"2015-6-7 21:23:8","gameOrder":"5574319c1d4bd4280909693f:c67d331f-62ba-466f-b2c9-b9f96660c195","sign":"34e5136a6f68e0e8dd2f9a84823605c4"}
I0607 21:26:05.299764 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"56266","order":"ed4a576c-757b-554b-a","price":10,"payType":1,"payCode":100013,"state":"success","time":"2015-6-7 21:25:40","gameOrder":"5574319c1d4bd4280909693f:faae3bca-dde3-48d4-ac9c-e3fc68be07dd","sign":"d0d2209492f12f4c83dd9b66ba3fbe56"}
I0607 21:28:41.799591 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"56266","order":"f97cbdf3-17f8-832d-4","price":10,"payType":1,"payCode":100013,"state":"success","time":"2015-6-7 21:28:35","gameOrder":"5574319c1d4bd4280909693f:e501eb10-f9ed-4d26-81a4-d85fdeb7cd59","sign":"80aeec481ff03472cf175e7dc0107dd7"}
I0607 21:39:33.184534 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1828538","order":"7c2b0083-2c58-11cf-3","price":10,"payType":1,"payCode":100013,"state":"success","time":"2015-6-7 21:39:15","gameOrder":"556b22111d4bd42809081b61:d66e0e2a-e54a-4c20-911b-ce738a75d466","sign":"1471829795650573cab2db6787ae96c6"}
I0607 23:10:56.375270 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"701301","order":"00634724-1028-6303-d","price":10,"payType":2,"payCode":100013,"state":"success","time":"2015-6-7 23:8:40","gameOrder":"5555943a1d4bd4280903bbcc:e63966d9-d4b7-49fb-bc84-125953ede500","sign":"15b3b03a0d27688c0dc2b9b05ed52112"}
I0607 23:18:39.349877 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"701301","order":"43dec979-62a2-1008-9","price":10,"payType":2,"payCode":100013,"state":"success","time":"2015-6-7 23:18:13","gameOrder":"5555943a1d4bd4280903bbcc:6f605bce-e19c-479c-bebc-817cca051ed9","sign":"4295e93abc19c35749b29f2df367fb06"}
I0607 23:19:42.284739 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"701301","order":"0941911e-8994-ff69-9","price":10,"payType":2,"payCode":100013,"state":"success","time":"2015-6-7 23:19:14","gameOrder":"5555943a1d4bd4280903bbcc:683d5a33-f7dc-4fda-b92f-0bef9ca38faa","sign":"b7ad287229883d581ce0197b1feb387a"}
I0607 23:20:48.074192 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"701301","order":"c2f2e168-0450-dba0-d","price":10,"payType":2,"payCode":100013,"state":"success","time":"2015-6-7 23:20:23","gameOrder":"5555943a1d4bd4280903bbcc:fd698cb9-9fc9-4ec1-a448-8649bd6f6b3f","sign":"034f7aea362abec64f6a0e4c6d0be50a"}
I0607 23:21:08.777799 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1313550","order":"8d32a0e3-03a0-816d-f","price":6,"payType":1,"payCode":100012,"state":"success","time":"2015-6-7 23:20:50","gameOrder":"55727bde1d4bd42809092260:a0f4a7e5-fcb9-4a45-b41a-45fedb5e067e","sign":"17fb6c504cd2ff05cd2f4c7d6d7f7db7"}
I0607 23:24:12.231447 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"701301","order":"d4b0be88-e08d-cf69-1","price":10,"payType":2,"payCode":100013,"state":"success","time":"2015-6-7 23:23:37","gameOrder":"5555943a1d4bd4280903bbcc:2187505f-4306-4da9-9318-b6daaec64643","sign":"58bdbf3cc737a347dcb9b544dff50e37"}
I0607 23:24:42.025952 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"701301","order":"f7c5eaa9-ad83-b161-5","price":10,"payType":2,"payCode":100013,"state":"success","time":"2015-6-7 23:24:17","gameOrder":"5555943a1d4bd4280903bbcc:3a6e2982-3e79-4747-9ca4-2397c6d9abd7","sign":"14a0836fed0c747839f035a45a2c5f47"}
I0607 23:26:27.254909 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"701301","order":"296c5569-ee1d-4296-0","price":10,"payType":2,"payCode":100013,"state":"success","time":"2015-6-7 23:25:56","gameOrder":"5555943a1d4bd4280903bbcc:fe8edcff-dc10-453f-8730-197f91902fb2","sign":"300364e7976b38a65c3993173fe9f3bf"}
I0607 23:28:29.448906 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1164998","order":"6ff83aca-3639-474c-0","price":10,"payType":1,"payCode":100013,"state":"success","time":"2015-6-7 23:28:7","gameOrder":"5569c9e11d4bd4280907ea31:563f8e84-6815-466f-9db9-3b9c682d4e67","sign":"096e8e86ef997935ded6f3e0957a76b0"}
I0607 23:30:56.629521 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1164998","order":"b75d35ba-3ca1-eab0-1","price":10,"payType":1,"payCode":100013,"state":"success","time":"2015-6-7 23:30:31","gameOrder":"5569c9e11d4bd4280907ea31:e06cd04a-aacf-4371-9685-63fc7e2f9d30","sign":"2ca80fa3103ccc38e71dc7a789528306"}
I0608 01:39:58.317524 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1098425","order":"88975214-b587-64fa-a","price":10,"payType":1,"payCode":100013,"state":"success","time":"2015-6-8 1:39:18","gameOrder":"555353101d4bd428090360db:cce4633d-0840-4aa2-9764-dbc1f419aea7","sign":"fdb7fd873b77aa276e1de1265c7d95b5"}
I0608 01:42:56.747816 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1098425","order":"8372bb27-9b23-d839-a","price":10,"payType":1,"payCode":100013,"state":"success","time":"2015-6-8 1:42:43","gameOrder":"555353101d4bd428090360db:076e3f0e-35c8-4d42-b08f-c27d34dc5c4c","sign":"32c9701ef8e4d044ceb0361e5eaf478f"}
I0608 03:17:26.343310 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1164998","order":"f7bdad62-4bfb-33c5-7","price":10,"payType":1,"payCode":100013,"state":"success","time":"2015-6-8 3:17:9","gameOrder":"5569c9e11d4bd4280907ea31:20eb4a67-0d0d-4133-b804-8b1eb0a91b87","sign":"d829c7c8d2060146170502d3101286a2"}
I0608 04:20:44.922642 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1822913","order":"b0c0f770-7e6f-ce14-0","price":10,"payType":1,"payCode":100013,"state":"success","time":"2015-6-8 4:20:29","gameOrder":"5564b4c01d4bd42809072e99:ab5ff01c-6944-4157-9ec1-f844752e6e2f","sign":"b244f661793bebcc656d1a251cb75826"}
I0608 08:45:43.745549 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"170254","order":"300ee48f-5b0e-9dc0-c","price":10,"payType":3,"payCode":100013,"state":"success","time":"2015-6-8 8:40:20","gameOrder":"551c16f11d4bd434e500687d:53fb1060-3254-406d-bde6-18208be4acf6","sign":"cb11101a6cb6e1fa2acb20dadb3dd995"}
I0608 09:09:41.622031 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1164998","order":"1d88b8b6-8249-d341-d","price":10,"payType":1,"payCode":100013,"state":"success","time":"2015-6-8 9:9:13","gameOrder":"5569c9e11d4bd4280907ea31:6ce982bb-12c6-4233-b3e3-18dd9aef5174","sign":"a6c91bf46eba60ca4d652457cf8ae318"}
I0608 09:11:17.574743 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1164998","order":"83143398-4a64-07d7-b","price":10,"payType":1,"payCode":100013,"state":"success","time":"2015-6-8 9:10:55","gameOrder":"5569c9e11d4bd4280907ea31:60d861b6-b55e-4a24-bbd4-e7bc6c092a08","sign":"1f4027700b26ea0853f29412bb6e5300"}
I0608 09:12:02.193387 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1164998","order":"450511aa-eaab-d86e-8","price":10,"payType":1,"payCode":100013,"state":"success","time":"2015-6-8 9:11:42","gameOrder":"5569c9e11d4bd4280907ea31:57b3e18c-15dc-4329-84cf-05b0cda5a651","sign":"7e516f1486d6a848851ef73571ddec94"}
I0608 11:21:36.268623 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1221062","order":"0dbad0b4-b547-8c56-5","price":10,"payType":1,"payCode":100013,"state":"success","time":"2015-6-8 11:21:18","gameOrder":"551b33be1d4bd434e5003629:9e9c95c0-9ca4-46d4-8a58-edfdb4022444","sign":"c2097aa81ad480b006b483a2e8985a64"}
I0608 11:24:13.692598 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1221062","order":"00ad2462-8c40-9217-b","price":10,"payType":1,"payCode":100013,"state":"success","time":"2015-6-8 11:23:0","gameOrder":"551b33be1d4bd434e5003629:ecabcd04-c7f5-43b6-99b0-cc617249093d","sign":"f78a0b939c8744bf9e45c5da1c288457"}
I0608 11:26:04.234195 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1221062","order":"afc92ee2-3794-e0c8-1","price":10,"payType":1,"payCode":100013,"state":"success","time":"2015-6-8 11:25:54","gameOrder":"551b33be1d4bd434e5003629:557e0d36-bd00-4233-8059-a9eac542ac59","sign":"ec524432d2c4e3d42bc7545ad5c8cf74"}
I0608 12:41:58.465228 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1684724","order":"1da54417-5564-01e7-0","price":6,"payType":1,"payCode":100012,"state":"success","time":"2015-6-8 12:41:32","gameOrder":"5552c1be1d4bd42809034227:4651c728-811f-4b18-8447-924a6548ba6c","sign":"3e1fbadef7130018339df8873b71120c"}
I0608 12:42:39.915257 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1684724","order":"d9879835-fe3a-ec71-6","price":6,"payType":1,"payCode":100012,"state":"success","time":"2015-6-8 12:42:14","gameOrder":"5552c1be1d4bd42809034227:83b1cf2b-54f1-47fa-9753-2c8277f7fc20","sign":"1a815b7bf4b5bdf64dcf74a1b844cabb"}
I0608 12:43:06.002612 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1684724","order":"bf8cdcf6-c805-2d96-1","price":6,"payType":1,"payCode":100012,"state":"success","time":"2015-6-8 12:42:47","gameOrder":"5552c1be1d4bd42809034227:803aa132-f430-4cf1-9238-469b2a1e3761","sign":"37a2103a7ce572f59f00856bfc075e22"}
I0608 12:44:16.418291 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1684724","order":"e00f31ab-cad2-2fe3-d","price":6,"payType":1,"payCode":100012,"state":"success","time":"2015-6-8 12:43:49","gameOrder":"5552c1be1d4bd42809034227:7734a302-7aeb-4b8c-a244-3f6b53779e2e","sign":"7d43c26d7fd6893f832b60726e81f235"}
I0608 12:45:16.635299 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1467152","order":"d59e56e6-b8e2-1e55-9","price":10,"payType":1,"payCode":100013,"state":"success","time":"2015-6-8 12:45:5","gameOrder":"555093a31d4bd4280902f259:81449a3f-7777-4996-aa37-ef165048bc2a","sign":"0096b5deadff3ffb84b23d124f9935be"}
I0608 13:06:50.995949 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1549479","order":"98756576-2e3f-56e1-2","price":10,"payType":1,"payCode":100013,"state":"success","time":"2015-6-8 13:5:4","gameOrder":"554416ac1d4bd4280900ba00:5ca781af-c000-4bed-b91d-f4ace5e1e1ce","sign":"d0850fb4037d2887eca5d13c1ae6f7ee"}
I0608 13:26:29.171505 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1765816","order":"fa4c5d1d-82a5-3ce1-5","price":6,"payType":1,"payCode":100012,"state":"success","time":"2015-6-8 13:26:5","gameOrder":"557524481d4bd42809098890:2a7be6e7-b337-4586-a73e-880353e3f5c3","sign":"a39f3b3190bb4a198074197b074f3d8c"}
I0608 14:17:47.168337 10249 qf_pay_handler.go:55] 起凡充值:{"appId":10061,"userId":"1765816","order":"d6e67685-bba6-eed0-d","price":6,"payType":1,"payCode":100012,"state":"success","time":"2015-6-8 14:12:11","gameOrder":"557524481d4bd42809098890:7ca40cd1-53b4-48e3-9256-b6d96f33e095","sign":"33bfef3e59186748ce94697019d130e9"}
`

type QfRes struct {
	AppId     interface{}      `json:"appId"`
	UserId    interface{} `json:"userId"`
	Order     interface{}      `json:"order"`
	Price     interface{} `json:"price"`
	PayType   interface{} `json:"payType"`
	PayCode   interface{} `json:"payCode"`
	State     interface{}      `json:"state"`
	Time      interface{}      `json:"time"`
	GameOrder interface{}      `json:"gameOrder"`
	Sign      interface{}      `json:"sign"`
}

type PrizeMail struct {
	UserId    string `json:"userId"`
	Gold      int    `json:"gold"`
	Diamond   int    `json:"diamond"`
	Exp       int    `json:"exp"`
	Score     int    `json:"score"`
	ItemType  int    `json:"itemType"`
	ItemCount int    `json:"itemCount"`
	Content   string `json:"content"`
}

func sendMail(userId string, diamond int) {
	mail := &PrizeMail{}
	mail.UserId = userId
	mail.Gold =  500000
	mail.Diamond = diamond
	mail.Exp = 0
	mail.Score = 0
	mail.Content = "起凡用户充值未到账补偿"
	b, err := json.Marshal(mail)
	if err != nil {
		fmt.Println("邮件打包失败err:", err, " mail:", mail)
		return
	}
	data := strings.NewReader(string(b))

	req, err := http.NewRequest("POST", "http://10.232.65.90:8002/sendPrizeMail?key=f1b9df1ed816c76ecfb2acf1c65b2a0d", data)
	if err != nil {
		fmt.Println("NewRequest失败:", err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Client.Do failed err:", err)
		return
	}
	io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()
}

func main() {
	flag.Parse()
	re := regexp.MustCompile(`{.*}`)
	res := re.FindAllString(s, -1);
	for _, item := range res {
		res := &QfRes{}
		err := json.Unmarshal([]byte(item), res)
		if err != nil {
			fmt.Println("解析失败err:", err, " item:", item)
			continue
		}
		u, err := user.FindByUserName(fmt.Sprintf("%v", res.UserId))
		if err != nil {
			fmt.Println("找不到玩家:", err, " userName:", res.UserId)
			continue
		}
		price := fmt.Sprintf("%v", res.Price)
		amount, err := strconv.ParseFloat(price, 64)
		if err != nil {
			continue
		}
		sendMail(u.UserId, int(amount))
		fmt.Println("用户Id:", u.UserId, " 充值金额:", res.Price, " 时间:", res.Time)
	}
}
