(window.webpackJsonp=window.webpackJsonp||[]).push([[5],{335:function(e,t,r){e.exports=r.p+"img/point.77f0338.png"},345:function(e,t,r){e.exports=r.p+"img/grid.43c7c41.png"},346:function(e,t,r){e.exports=r.p+"img/nkn-logo.8e7be89.png"},347:function(e,t,r){"use strict";r.d(t,"a",(function(){return n}));r(22);function n(e){for(var t=arguments.length,r=new Array(t>1?t-1:0),n=1;n<t;n++)r[n-1]=arguments[n];for(var o=0,c=r;o<c.length;o++){var source=c[o];if(source)for(var l=0,d=Object.keys(source);l<d.length;l++){var v=d[l];void 0!==source[v]&&(e[v]=source[v])}}return e}},377:function(e,t,r){var content=r(410);content.__esModule&&(content=content.default),"string"==typeof content&&(content=[[e.i,content,""]]),content.locals&&(e.exports=content.locals);(0,r(52).default)("04288992",content,!0,{sourceMap:!1})},409:function(e,t,r){"use strict";r(377)},410:function(e,t,r){var n=r(51)(!1);n.push([e.i,".v-btn{text-transform:none!important}.v-text-field.v-text-field--solo:not(.v-text-field--solo-flat)>.v-input__control>.v-input__slot{border:1.08px solid!important;border-image-source:linear-gradient(.606turn,hsla(0,0%,100%,.3) -27.92%,hsla(0,0%,100%,0) 92.11%)!important;background:hsla(0,0%,100%,.15)!important;box-shadow:0 17.5609px 16.2101px rgba(11,39,40,.07)!important;-webkit-backdrop-filter:blur(18.9118px)!important;backdrop-filter:blur(18.9118px)!important;border-radius:10px!important}.v-text-field.v-text-field--solo .v-input__append-outer,.v-text-field.v-text-field--solo .v-input__prepend-outer{margin-top:6px!important}.v-text-field.v-text-field--solo .v-input__control{min-height:36px!important}.v-select__selection--comma{text-shadow:0 4px 4px rgba(0,0,0,.25)!important}header{height:140px;display:flex;align-items:center}.container *{z-index:1!important}.nav .nav-item{margin-bottom:10px!important;margin-left:10px!important}table{width:100%}table,td,th{margin-top:1rem;border:1px solid #888;border-collapse:collapse;text-align:center;padding-left:10px;padding-right:10px;font-size:small}td,th{width:30px}tr{height:3rem}tr:hover{background-color:#555}.expand{width:4rem;text-align:center}.expand:hover{background-color:#555;cursor:pointer}.blankline{display:block;height:1rem}",""]),e.exports=n},522:function(e,t,r){"use strict";r.r(t);var n=r(16),o=(r(6),r(69),r(68),r(70),r(87),r(101)),c=r.n(o),l=r(347),d="/rpc/network",v={getNetworkConfig:{method:"getNetworkConfig"},setNetworkConfig:{method:"setNetworkConfig"},authorizeMember:{method:"authorizeMember"},removeMember:{method:"removeMember"},deleteWaiting:{method:"deleteWaiting"},setAcceptAddress:{method:"setAcceptAddress"},sendToken:{method:"sendToken"},nknPing:{method:"nknPing"}},m={},f=function(e){v.hasOwnProperty(e)&&(m[e]=function(t,r){return r=l.a({},v[e].defaultParams,r),function(e,t){return w.apply(this,arguments)}(t,v[e].method,r)})};for(var h in v)f(h);function w(){return w=Object(n.a)(regeneratorRuntime.mark((function e(t,r){var n,o,l,data,d=arguments;return regeneratorRuntime.wrap((function(e){for(;;)switch(e.prev=e.next){case 0:return n=d.length>2&&void 0!==d[2]?d[2]:{},e.prev=1,e.next=4,window.rpcHeaders;case 4:o=e.sent,e.next=10;break;case 7:e.prev=7,e.t0=e.catch(1),console.error("Await rpc headers error:",e.t0);case 10:return e.next=12,c()({url:t,method:"POST",timeout:1e4,headers:o,data:{id:"nConnect-web",jsonrpc:"2.0",method:r,params:n}});case 12:if(l=e.sent,!(data=l.data).error){e.next=16;break}throw data.error;case 16:if(void 0===data.result){e.next=18;break}return e.abrupt("return",data.result);case 18:throw new Error("rpc response contains no result or error field");case 19:case"end":return e.stop()}}),e,null,[[1,7]])}))),w.apply(this,arguments)}function k(){return _.apply(this,arguments)}function _(){return(_=Object(n.a)(regeneratorRuntime.mark((function e(){return regeneratorRuntime.wrap((function(e){for(;;)switch(e.prev=e.next){case 0:return e.abrupt("return",m.getNetworkConfig(d));case 1:case"end":return e.stop()}}),e)})))).apply(this,arguments)}function x(e){return A.apply(this,arguments)}function A(){return(A=Object(n.a)(regeneratorRuntime.mark((function e(t){return regeneratorRuntime.wrap((function(e){for(;;)switch(e.prev=e.next){case 0:return e.abrupt("return",m.setNetworkConfig(d,{domain:t.domain,ipStart:t.ipStart,ipEnd:t.ipEnd,netmask:t.netmask,gateway:t.gateway,dns:t.dns}));case 1:case"end":return e.stop()}}),e)})))).apply(this,arguments)}function y(e){return C.apply(this,arguments)}function C(){return(C=Object(n.a)(regeneratorRuntime.mark((function e(address){return regeneratorRuntime.wrap((function(e){for(;;)switch(e.prev=e.next){case 0:return e.abrupt("return",m.authorizeMember(d,{address:address}));case 1:case"end":return e.stop()}}),e)})))).apply(this,arguments)}function R(e){return D.apply(this,arguments)}function D(){return(D=Object(n.a)(regeneratorRuntime.mark((function e(address){return regeneratorRuntime.wrap((function(e){for(;;)switch(e.prev=e.next){case 0:return e.abrupt("return",m.removeMember(d,{address:address}));case 1:case"end":return e.stop()}}),e)})))).apply(this,arguments)}function S(e){return j.apply(this,arguments)}function j(){return(j=Object(n.a)(regeneratorRuntime.mark((function e(address){return regeneratorRuntime.wrap((function(e){for(;;)switch(e.prev=e.next){case 0:return e.abrupt("return",m.deleteWaiting(d,{address:address}));case 1:case"end":return e.stop()}}),e)})))).apply(this,arguments)}function N(e,t){return O.apply(this,arguments)}function O(){return(O=Object(n.a)(regeneratorRuntime.mark((function e(address,t){return regeneratorRuntime.wrap((function(e){for(;;)switch(e.prev=e.next){case 0:return e.abrupt("return",m.setAcceptAddress(d,{address:address,AcceptAddresses:t}));case 1:case"end":return e.stop()}}),e)})))).apply(this,arguments)}function M(e,t){return I.apply(this,arguments)}function I(){return(I=Object(n.a)(regeneratorRuntime.mark((function e(address,t){return regeneratorRuntime.wrap((function(e){for(;;)switch(e.prev=e.next){case 0:return e.abrupt("return",m.sendToken(d,{address:address,amount:t}));case 1:case"end":return e.stop()}}),e)})))).apply(this,arguments)}function P(e){return V.apply(this,arguments)}function V(){return(V=Object(n.a)(regeneratorRuntime.mark((function e(address){return regeneratorRuntime.wrap((function(e){for(;;)switch(e.prev=e.next){case 0:return e.abrupt("return",m.nknPing(d,{address:address}));case 1:case"end":return e.stop()}}),e)})))).apply(this,arguments)}var T=r(167),$=r.n(T),z={name:"network",computed:{availableLocales:function(){return this.$i18n.locales}},data:function(){return{lang:this.$i18n.locale,managerAddress:"",managerBalance:"",showConfig:!1,showMembers:!0,networkData:{networkInfo:{},member:[{name:"bill",ip:"10.0.86.3",address:"aaaabbbbccccdddd",lastSeen:"2023-09-11 13:00:00"}]},waitingCheck:[],memberCheck:[],setAddress:"",acceptAddress:[],allAddress:[],dialog:!1,showResponse:!1,response:"",confirmDialog:!1,confirmTitle:"",confirm:!1,sendDialog:!1,amount:1}},mounted:function(){var e=this;return Object(n.a)(regeneratorRuntime.mark((function t(){return regeneratorRuntime.wrap((function(t){for(;;)switch(t.prev=t.next){case 0:e.getNetworkConfig();case 1:case"end":return t.stop()}}),t)})))()},created:function(){return Object(n.a)(regeneratorRuntime.mark((function e(){return regeneratorRuntime.wrap((function(e){for(;;)switch(e.prev=e.next){case 0:case"end":return e.stop()}}),e)})))()},methods:{onChangeSwitchLanguage:function(e){this.$i18n.locale=e,$.a.set("language",e)},getNetworkConfig:function(){var e=this;return Object(n.a)(regeneratorRuntime.mark((function t(){var r;return regeneratorRuntime.wrap((function(t){for(;;)switch(t.prev=t.next){case 0:return t.prev=0,t.next=3,k();case 3:r=t.sent,e.networkData=r.networkData,e.managerAddress=r.managerAddress,e.managerBalance=r.managerBalance,t.next=13;break;case 9:t.prev=9,t.t0=t.catch(0),console.error(t.t0),window.alert(t.t0);case 13:case"end":return t.stop()}}),t,null,[[0,9]])})))()},setNetworkConfig:function(){var e=this;return Object(n.a)(regeneratorRuntime.mark((function t(){var r;return regeneratorRuntime.wrap((function(t){for(;;)switch(t.prev=t.next){case 0:return t.prev=0,t.next=3,x(e.networkData);case 3:r=t.sent,e.response=r,e.showResponse=!0,t.next=12;break;case 8:t.prev=8,t.t0=t.catch(0),console.error(t.t0),window.alert(t.t0);case 12:case"end":return t.stop()}}),t,null,[[0,8]])})))()},authorizeMember:function(address){var e=this;return Object(n.a)(regeneratorRuntime.mark((function t(){return regeneratorRuntime.wrap((function(t){for(;;)switch(t.prev=t.next){case 0:return t.next=2,y(address);case 2:"success"==t.sent&&(e.getNetworkConfig(),e.toSetAcceptAddress(address));case 4:case"end":return t.stop()}}),t)})))()},removeMember:function(address){var e=this;return Object(n.a)(regeneratorRuntime.mark((function t(){var r;return regeneratorRuntime.wrap((function(t){for(;;)switch(t.prev=t.next){case 0:return(r=e).confirm=!1,r.confirmTitle="Are you sure to remove this member?",r.confirmDialog=!0,t.next=6,new Promise((function(e,t){var n=setInterval((function(){0==r.confirmDialog&&(clearInterval(n),r.confirm?e():t())}),100)})).then(Object(n.a)(regeneratorRuntime.mark((function e(){var t,n;return regeneratorRuntime.wrap((function(e){for(;;)switch(e.prev=e.next){case 0:return e.next=2,R(address);case 2:"success"==(t=e.sent)&&(n=r.networkData.member[address],delete r.networkData.member[address],r.networkData.waiting[address]=n,r.response=t,r.showResponse=!0);case 4:case"end":return e.stop()}}),e)})))).catch((function(){console.log("cancel")}));case 6:case"end":return t.stop()}}),t)})))()},sendToken:function(address){var e=this;return Object(n.a)(regeneratorRuntime.mark((function t(){var r;return regeneratorRuntime.wrap((function(t){for(;;)switch(t.prev=t.next){case 0:return(r=e).confirm=!1,r.sendDialog=!0,t.next=5,new Promise((function(e,t){var n=setInterval((function(){0==r.sendDialog&&(clearInterval(n),console.log("confirm",r.confirm),r.confirm?e():t())}),100)})).then(Object(n.a)(regeneratorRuntime.mark((function t(){var n;return regeneratorRuntime.wrap((function(t){for(;;)switch(t.prev=t.next){case 0:if(!(e.amount>0)){t.next=7;break}return console.log("send",e.amount," token to ",address),t.next=4,M(address,e.amount);case 4:n=t.sent,r.response=n,r.showResponse=!0;case 7:case"end":return t.stop()}}),t)})))).catch((function(){console.log("cancel")}));case 5:case"end":return t.stop()}}),t)})))()},deleteWaiting:function(address){var e=this;return Object(n.a)(regeneratorRuntime.mark((function t(){var r;return regeneratorRuntime.wrap((function(t){for(;;)switch(t.prev=t.next){case 0:return t.next=2,S(address);case 2:"success"==(r=t.sent)&&(delete e.networkData.waiting[address],e.response=r,e.showResponse=!0);case 4:case"end":return t.stop()}}),t)})))()},nknPing:function(address){var e=this;return Object(n.a)(regeneratorRuntime.mark((function t(){var r;return regeneratorRuntime.wrap((function(t){for(;;)switch(t.prev=t.next){case 0:return t.next=2,P(address);case 2:(r=t.sent).includes("success")&&(e.response=r,e.showResponse=!0);case 4:case"end":return t.stop()}}),t)})))()},clickAllAddress:function(){if(this.allAddress.length>0)for(var address in this.acceptAddress=[],this.networkData.member)this.acceptAddress.push(address);else this.acceptAddress=[]},selectAddress:function(address){this.acceptAddress.includes(address)||(this.allAddress=[])},toSetAcceptAddress:function(address){if(this.setAddress=address,this.allAddress=[],this.acceptAddress=this.networkData.acceptAddress[address],this.acceptAddress||(this.acceptAddress=[]),this.acceptAddress&&this.acceptAddress.length>0&&"allMembers"==this.acceptAddress[0])for(var e in this.allAddress=["allMembers"],this.acceptAddress=[],this.networkData.member)this.acceptAddress.push(e);this.dialog=!0},setAcceptAddress:function(){var e=this;return Object(n.a)(regeneratorRuntime.mark((function t(){var r,n;return regeneratorRuntime.wrap((function(t){for(;;)switch(t.prev=t.next){case 0:return e.dialog=!1,r=e.acceptAddress,e.allAddress.length>0&&(r=["allMembers"]),t.next=5,N(e.setAddress,r);case 5:"success"==(n=t.sent)&&(e.networkData.acceptAddress[e.setAddress]=r),e.response=n,e.showResponse=!0;case 9:case"end":return t.stop()}}),t)})))()}}},E=(r(409),r(64)),B=r(77),L=r.n(B),W=r(525),F=r(382),K=r(357),J=r(526),G=r(514),H=r(327),Y=r(524),Q=r(391),U=r(527),X=r(528),Z=r(516),ee=r(521),te=r(517),re=r(401),component=Object(E.a)(z,(function(){var e=this,t=e.$createElement,n=e._self._c||t;return n("v-container",[n("v-parallax",{staticStyle:{position:"absolute",top:"0","max-width":"1600px","min-width":"640px",width:"80%"},attrs:{src:r(345),height:"800"}},[n("v-img",{staticStyle:{position:"absolute",right:"38%",bottom:"130px"},attrs:{src:r(335),width:"40",height:"40",contain:""}}),e._v(" "),n("v-img",{staticStyle:{position:"absolute",right:"20%",bottom:"80px"},attrs:{src:r(335),width:"67",height:"67",contain:""}})],1),e._v(" "),n("v-container",{staticClass:"container"},[n("header",[n("v-row",{attrs:{align:"center"}},[n("v-col",[n("img",{attrs:{src:r(346),alt:""}})]),e._v(" "),n("v-col",{staticClass:"ml-auto",attrs:{cols:"auto"}},[n("span",[e._v("Language")]),e._v(" "),n("span",{staticStyle:{"max-width":"200px",display:"inline-flex"}},[n("v-select",{staticClass:"nav-item select-language",staticStyle:{width:"200px"},attrs:{items:e.availableLocales,"item-value":"code","item-text":"name",label:"Language",solo:"",dark:""},on:{change:function(t){return e.onChangeSwitchLanguage(e.lang)}},model:{value:e.lang,callback:function(t){e.lang=t},expression:"lang"}})],1)])],1)],1),e._v(" "),n("v-row",[n("v-col",{attrs:{cols:"20"}},[n("div",{staticClass:"bg-linear-1 pa-4 mb-4"},[n("div",{staticClass:"w-100 d-flex justify-space-between mb-4"},[n("h3",[e._v("Network Manager")])]),e._v(" "),n("div",{staticClass:"d-flex flex-wrap justify-space-between mt-5"},[n("div",[e._v("Address: "+e._s(e.managerAddress))]),e._v(" "),n("div",[e._v("NKN Balance: "+e._s(Number(e.managerBalance).toFixed(2)))])])]),e._v(" "),n("div",{staticClass:"bg-linear-1 pa-4 mb-4"},[n("div",{staticClass:"w-100 d-flex justify-space-between mb-4"},[n("h3",[e._v("Network Configuration")]),e._v(" "),n("div",{staticClass:"expand",on:{click:function(t){e.showConfig=!e.showConfig}}},[e._v("▼")])]),e._v(" "),e.showConfig?n("div",[n("v-text-field",{staticClass:"mx-6",attrs:{clearable:"","hide-details":"auto",label:"Network domain"},model:{value:e.networkData.networkInfo.domain,callback:function(t){e.$set(e.networkData.networkInfo,"domain",t)},expression:"networkData.networkInfo.domain"}}),e._v(" "),n("div",{staticClass:"d-flex flex-wrap mt-5"},[n("v-text-field",{staticClass:"mx-6",attrs:{clearable:"","hide-details":"auto",label:"IP Start"},model:{value:e.networkData.ipStart,callback:function(t){e.$set(e.networkData,"ipStart",t)},expression:"networkData.ipStart"}}),e._v(" "),n("v-text-field",{staticClass:"mx-6",attrs:{clearable:"","hide-details":"auto",label:"IP End"},model:{value:e.networkData.ipEnd,callback:function(t){e.$set(e.networkData,"ipEnd",t)},expression:"networkData.ipEnd"}}),e._v(" "),n("v-text-field",{staticClass:"mx-6",attrs:{clearable:"","hide-details":"auto",label:"Network mask"},model:{value:e.networkData.netmask,callback:function(t){e.$set(e.networkData,"netmask",t)},expression:"networkData.netmask"}})],1),e._v(" "),n("div",{staticClass:"d-flex flex-wrap mt-5"},[n("v-text-field",{staticClass:"mx-6",attrs:{clearable:"","hide-details":"auto",label:"Gateway"},model:{value:e.networkData.networkInfo.gateway,callback:function(t){e.$set(e.networkData.networkInfo,"gateway",t)},expression:"networkData.networkInfo.gateway"}}),e._v(" "),n("v-text-field",{staticClass:"mx-6",attrs:{clearable:"","hide-details":"auto",label:"DNS"},model:{value:e.networkData.networkInfo.dns,callback:function(t){e.$set(e.networkData.networkInfo,"dns",t)},expression:"networkData.networkInfo.dns"}})],1),e._v(" "),n("div",{staticClass:"d-flex justify-center mt-12"},[n("v-btn",{staticStyle:{width:"40%"},attrs:{color:"primary"},on:{click:e.setNetworkConfig}},[e._v("Submit")])],1)],1):e._e()]),e._v(" "),n("div",{staticClass:"bg-linear-1 pa-4 mb-4",staticStyle:{"overflow-x":"scroll"}},[n("div",{staticClass:"w-100 d-flex justify-space-between mb-4"},[n("h3",[e._v("Waiting for Authorization")]),e._v(" "),n("v-btn",{attrs:{color:"primary"},on:{click:e.getNetworkConfig}},[e._v("Refresh")])],1),e._v(" "),n("table",[n("tr",[n("th",[e._v("Name")]),e._v(" "),n("th",[e._v("Address")]),e._v(" "),n("th",[e._v("Accept")]),e._v(" "),n("th",[e._v("Reject")])]),e._v(" "),e._l(e.networkData.waiting,(function(t){return n("tr",{key:t.address},[n("td",[e._v(e._s(t.name))]),e._v(" "),n("td",[e._v(e._s(t.address))]),e._v(" "),n("td",[n("v-btn",{staticClass:"ma-2",attrs:{color:"primary"},on:{click:function(r){return e.authorizeMember(t.address)}}},[e._v(" Accept ")])],1),e._v(" "),n("td",[n("v-btn",{attrs:{color:"secondary"},on:{click:function(r){return e.deleteWaiting(t.address)}}},[e._v("Reject")])],1)])}))],2)]),e._v(" "),n("div",{staticClass:"bg-linear-1 pa-4 mb-4",staticStyle:{"overflow-x":"scroll"}},[n("div",{staticClass:"w-100 d-flex justify-space-between mb-4"},[n("h3",[e._v("Network Members")]),e._v(" "),n("div",{staticClass:"expand",on:{click:function(t){e.showMembers=!e.showMembers}}},[e._v("▼")])]),e._v(" "),e.showMembers?n("table",[n("tr",[n("th",[e._v("Name")]),e._v(" "),n("th",[e._v("IP")]),e._v(" "),n("th",[e._v("LastSeen")]),e._v(" "),n("th",[e._v("Server")]),e._v(" "),n("th",[e._v("Balance")]),e._v(" "),n("th",[e._v("Address")]),e._v(" "),n("th",[e._v("Accept")]),e._v(" "),n("th",[e._v("Send Token")]),e._v(" "),n("th",[e._v("Ping")]),e._v(" "),n("th",[e._v("Remove")])]),e._v(" "),e._l(e.networkData.member,(function(t){return n("tr",{key:t.address},[n("td",[e._v(e._s(t.name))]),e._v(" "),n("td",[e._v(e._s(t.ip))]),e._v(" "),n("td",[e._v(e._s(t.lastSeen.substring(2,19).replace("T"," ")))]),e._v(" "),n("td",{style:t.server?"background:green":""},[e._v(e._s(t.server?"Yes":"No"))]),e._v(" "),n("td",{style:t.server?t.balance>.1?"background:green":"background:orange":""},[e._v("\n                  "+e._s(Number(t.balance)>0?Number(t.balance).toFixed(2):t.balance))]),e._v(" "),n("td",{staticStyle:{width:"260px","word-break":"break-all"}},[e._v(e._s(t.address))]),e._v(" "),n("td",[n("v-btn",{attrs:{color:"primary"},on:{click:function(r){return e.toSetAcceptAddress(t.address)}}},[e._v("Set")])],1),e._v(" "),n("td",[t.server?n("v-btn",{attrs:{color:"primary"},on:{click:function(r){return e.sendToken(t.address)}}},[e._v("Send")]):e._e()],1),e._v(" "),n("td",[n("v-btn",{attrs:{color:"primary"},on:{click:function(r){return e.nknPing(t.address)}}},[e._v("Ping")])],1),e._v(" "),n("td",[n("v-btn",{attrs:{color:"secondary"},on:{click:function(r){return e.removeMember(t.address)}}},[e._v("Remove")])],1)])}))],2):e._e()])])],1),e._v(" "),n("v-dialog",{attrs:{width:"auto"},model:{value:e.dialog,callback:function(t){e.dialog=t},expression:"dialog"}},[n("v-card",[n("v-card-text",[n("h3",{staticStyle:{"margin-top":"1rem"}},[e._v("Set Accept Address")]),e._v(" "),n("h5",{staticClass:"mt-5"},[e._v("Accept All")]),e._v(" "),n("v-checkbox",{attrs:{label:"All members",color:"red",value:"allMembers"},on:{click:e.clickAllAddress},model:{value:e.allAddress,callback:function(t){e.allAddress=t},expression:"allAddress"}}),e._v(" "),n("v-divider"),e._v(" "),n("div",[n("table",[n("tr",[n("th",[e._v("Name")]),e._v(" "),n("th",[e._v("IP")]),e._v(" "),n("th",[e._v("Address")])]),e._v(" "),e._l(e.networkData.member,(function(t){return t.address!=e.setAddress?n("tr",{key:t.address},[n("td",[n("v-checkbox",{attrs:{label:t.name,color:"red",value:t.address},on:{click:function(r){return e.selectAddress(t.address)}},model:{value:e.acceptAddress,callback:function(t){e.acceptAddress=t},expression:"acceptAddress"}})],1),e._v(" "),n("td",[e._v(e._s(t.ip))]),e._v(" "),n("td",[e._v(e._s(t.address))])]):e._e()}))],2)])],1),e._v(" "),n("v-card-actions",[n("v-spacer"),e._v(" "),n("v-btn",{attrs:{color:"secondary"},on:{click:function(t){e.dialog=!1}}},[e._v(" Cancel ")]),e._v(" "),n("v-btn",{attrs:{color:"primary"},on:{click:e.setAcceptAddress}},[e._v(" Save ")])],1)],1)],1),e._v(" "),n("v-dialog",{attrs:{transition:"dialog-top-transition",width:"400"},model:{value:e.showResponse,callback:function(t){e.showResponse=t},expression:"showResponse"}},[n("v-card",{staticClass:"mx-auto"},[n("div",{staticClass:"blankline"}),e._v(" "),n("v-card-text",[n("h3",{staticClass:"mt-4 mb-4"},[e._v("Tips")]),e._v(" "),n("h3",[e._v(e._s(e.response))])]),e._v(" "),n("v-card-actions",[n("v-spacer"),e._v(" "),n("v-btn",{attrs:{color:"primary"},on:{click:function(t){e.showResponse=!1}}},[e._v(" Close ")])],1)],1)],1),e._v(" "),n("v-dialog",{attrs:{width:"300"},model:{value:e.confirmDialog,callback:function(t){e.confirmDialog=t},expression:"confirmDialog"}},[n("v-card",[n("v-card-text",[e._v(" "+e._s(e.confirmTitle)+" ")]),e._v(" "),n("v-card-actions",[n("v-spacer"),e._v(" "),n("v-btn",{attrs:{color:"secondary"},on:{click:function(t){e.confirmDialog=!1,e.confirm=!1}}},[e._v("Cancel")]),e._v(" "),n("v-btn",{attrs:{color:"primary"},on:{click:function(t){e.confirmDialog=!1,e.confirm=!0}}},[e._v("Confirm")])],1)],1)],1),e._v(" "),n("v-dialog",{attrs:{width:"300"},model:{value:e.sendDialog,callback:function(t){e.sendDialog=t},expression:"sendDialog"}},[n("v-card",[n("div",{staticClass:"blankline"}),e._v(" "),n("v-card-text",[e._v(" Send NKN token to member ")]),e._v(" "),n("div",{staticClass:"blankline"}),e._v(" "),n("v-text-field",{staticClass:"mx-6",attrs:{clearable:"","hide-details":"auto",label:"NKN Amount"},model:{value:e.amount,callback:function(t){e.amount=t},expression:"amount"}}),e._v(" "),n("div",{staticClass:"blankline"}),e._v(" "),n("v-card-actions",[n("v-btn",{attrs:{color:"secondary"},on:{click:function(t){e.sendDialog=!1,e.confirm=!1}}},[e._v("Cancel")]),e._v(" "),n("v-spacer"),e._v(" "),n("v-btn",{attrs:{color:"primary"},on:{click:function(t){e.sendDialog=!1,e.confirm=!0}}},[e._v("Confirm")])],1)],1)],1)],1)],1)}),[],!1,null,null,null);t.default=component.exports;L()(component,{VBtn:W.a,VCard:F.a,VCardActions:K.a,VCardText:K.b,VCheckbox:J.a,VCol:G.a,VContainer:H.a,VDialog:Y.a,VDivider:Q.a,VImg:U.a,VParallax:X.a,VRow:Z.a,VSelect:ee.a,VSpacer:te.a,VTextField:re.a})}}]);