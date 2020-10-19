(this["webpackJsonpnconnect-web"]=this["webpackJsonpnconnect-web"]||[]).push([[0],{123:function(e,t,n){},127:function(e,t,n){"use strict";n.r(t);var a=n(0),r=n.n(a),d=n(8),c=n.n(d),s=(n(78),n(79),n(10)),i=n.n(s),o=n(13),l=n(57),u=n(58),h=n(21),m=n(65),p=n(66),v=n(59),f=n.n(v),A=n(158),b=n(166),w=n(164),g=n(165),k=n(162),y=n(163),E=n(60),C=n.n(E);var S="/rpc/admin",j={getAdminToken:{method:"getAdminToken"},getAddrs:{method:"getAddrs"},setAddrs:{method:"setAddrs"},addAddrs:{method:"addAddrs"},removeAddrs:{method:"removeAddrs"},getLocalIP:{method:"getLocalIP"},getInfo:{method:"getInfo"}},O={},x=function(e){j.hasOwnProperty(e)&&(O[e]=function(t,n){return n=function(e){for(var t=arguments.length,n=new Array(t>1?t-1:0),a=1;a<t;a++)n[a-1]=arguments[a];for(var r=0,d=n;r<d.length;r++){var c=d[r];if(c)for(var s=0,i=Object.keys(c);s<i.length;s++){var o=i[s];void 0!==c[o]&&(e[o]=c[o])}}return e}({},j[e].defaultParams,n),function(e,t){return P.apply(this,arguments)}(t,j[e].method,n)})};for(var T in j)x(T);function P(){return(P=Object(o.a)(i.a.mark((function e(t,n){var a,r,d,c=arguments;return i.a.wrap((function(e){for(;;)switch(e.prev=e.next){case 0:return a=c.length>2&&void 0!==c[2]?c[2]:{},e.next=3,C()({url:t,method:"POST",timeout:1e4,data:{id:"nConnect-web",jsonrpc:"2.0",method:n,params:a}});case 3:if(r=e.sent,!(d=r.data).error){e.next=7;break}throw d.error;case 7:if(void 0===d.result){e.next=9;break}return e.abrupt("return",d.result);case 9:throw"rpc response contains no result or error field";case 10:case"end":return e.stop()}}),e)})))).apply(this,arguments)}function I(){return N.apply(this,arguments)}function N(){return(N=Object(o.a)(i.a.mark((function e(){return i.a.wrap((function(e){for(;;)switch(e.prev=e.next){case 0:return e.abrupt("return",O.getAdminToken(S));case 1:case"end":return e.stop()}}),e)})))).apply(this,arguments)}function R(){return(R=Object(o.a)(i.a.mark((function e(){return i.a.wrap((function(e){for(;;)switch(e.prev=e.next){case 0:return e.abrupt("return",O.getAddrs(S));case 1:case"end":return e.stop()}}),e)})))).apply(this,arguments)}function L(e,t){return M.apply(this,arguments)}function M(){return(M=Object(o.a)(i.a.mark((function e(t,n){var a;return i.a.wrap((function(e){for(;;)switch(e.prev=e.next){case 0:return a={},t&&(a.acceptAddrs=t),n&&(a.adminAddrs=n),e.abrupt("return",O.setAddrs(S,a));case 4:case"end":return e.stop()}}),e)})))).apply(this,arguments)}function Q(){return(Q=Object(o.a)(i.a.mark((function e(){return i.a.wrap((function(e){for(;;)switch(e.prev=e.next){case 0:return e.abrupt("return",O.getInfo(S));case 1:case"end":return e.stop()}}),e)})))).apply(this,arguments)}n(123);function D(e){return e?e.join("\n"):""}function J(e){return e?e.split("\n").filter((function(e){return e.length>0})):[]}function _(e){var t=e.split(".");return t[t.length-1]}var B=function(e){Object(p.a)(n,e);var t=Object(m.a)(n);function n(e){var a;return Object(l.a)(this,n),(a=t.call(this,e)).state={adminTokenStr:"",adminTokenQRCode:"",acceptAddrs:"",adminAddrs:"",addr:"",localIP:[],showAdvanced:!1},a.handleAcceptAddrsChange=a.handleAcceptAddrsChange.bind(Object(h.a)(a)),a.handleAdminAddrsChange=a.handleAdminAddrsChange.bind(Object(h.a)(a)),a.handleSubmit=a.handleSubmit.bind(Object(h.a)(a)),a.updateAdminToken=a.updateAdminToken.bind(Object(h.a)(a)),a.handleAdvancedChange=a.handleAdvancedChange.bind(Object(h.a)(a)),a}return Object(u.a)(n,[{key:"handleAcceptAddrsChange",value:function(e){this.setState({acceptAddrs:e.target.value})}},{key:"handleAdminAddrsChange",value:function(e){this.setState({adminAddrs:e.target.value})}},{key:"handleAdvancedChange",value:function(e){this.setState({showAdvanced:!this.state.showAdvanced}),console.log(this.state.showAdvanced)}},{key:"handleSubmit",value:function(){var e=Object(o.a)(i.a.mark((function e(t){var n;return i.a.wrap((function(e){for(;;)switch(e.prev=e.next){case 0:return t.preventDefault(),e.prev=1,e.next=4,L(J(this.state.acceptAddrs),J(this.state.adminAddrs));case 4:n=e.sent,this.setState({acceptAddrs:D(n.acceptAddrs),adminAddrs:D(n.adminAddrs)}),alert("Save success!"),e.next=13;break;case 9:e.prev=9,e.t0=e.catch(1),console.error(e.t0),alert(e.t0);case 13:case"end":return e.stop()}}),e,this,[[1,9]])})));return function(t){return e.apply(this,arguments)}}()},{key:"updateAdminToken",value:function(){var e=Object(o.a)(i.a.mark((function e(){var t,n,a;return i.a.wrap((function(e){for(;;)switch(e.prev=e.next){case 0:return e.prev=0,e.next=3,I();case 3:if(!(t=e.sent)){e.next=10;break}return n=JSON.stringify(t),e.next=8,f.a.toDataURL(n);case 8:a=e.sent,this.setState({adminTokenStr:n,adminTokenQRCode:a});case 10:e.next=15;break;case 12:e.prev=12,e.t0=e.catch(0),console.error(e.t0);case 15:case"end":return e.stop()}}),e,this,[[0,12]])})));return function(){return e.apply(this,arguments)}}()},{key:"componentDidMount",value:function(){var e=this;this.updateAdminToken(),setInterval(this.updateAdminToken,3e5),function(){return R.apply(this,arguments)}().then((function(t){e.setState({acceptAddrs:D(t.acceptAddrs),adminAddrs:D(t.adminAddrs)})})).catch((function(e){console.error(e),alert(e)})),function(){return Q.apply(this,arguments)}().then((function(t){e.setState({addr:t.addr,localIP:t.localIP.ipv4})})).catch((function(e){console.error(e)}))}},{key:"render",value:function(){return r.a.createElement("div",{className:"App"},r.a.createElement(A.a,null,r.a.createElement("div",{className:"row"},r.a.createElement("img",{src:"/static/media/nkn_logo.png"})),r.a.createElement("div",{className:"row"},r.a.createElement("img",{src:this.state.adminTokenQRCode})),r.a.createElement("div",{className:"row"},"Scan the QR code on nMobile Pro to connect and manage device."),r.a.createElement("div",{className:"row"},"Purchase data plan on nMobile Pro or ",r.a.createElement("a",{target:"_blank",href:"http://nconnect-payment.nkncdn.com/payment/?addr="+_(this.state.addr)},"web payment portal")),r.a.createElement("div",{className:"row"},r.a.createElement(b.a,{variant:"outlined",color:"primary",onClick:this.handleAdvancedChange,style:{width:"100%"}},this.state.showAdvanced?r.a.createElement(k.a,null):r.a.createElement(y.a,null),this.state.showAdvanced?"Hide Advanced":"Show Advanced")),r.a.createElement(w.a,{in:this.state.showAdvanced},r.a.createElement("div",{className:"advanced-row"},r.a.createElement(g.a,{disabled:!0,multiline:!0,label:"Local IP address",value:this.state.localIP.join("\n"),style:{width:"100%"}})),r.a.createElement("div",{className:"advanced-row"},r.a.createElement(g.a,{disabled:!0,multiline:!0,label:"Access key (expires in 5 minutes)",value:this.state.adminTokenStr,style:{width:"100%"}})),r.a.createElement("div",{className:"advanced-row"},r.a.createElement(g.a,{multiline:!0,variant:"filled",label:"Accept addresses",value:this.state.acceptAddrs,onChange:this.handleAcceptAddrsChange,style:{width:"100%"}}),r.a.createElement(g.a,{multiline:!0,variant:"filled",label:"Admins",value:this.state.adminAddrs,onChange:this.handleAdminAddrsChange,style:{width:"100%"}})),r.a.createElement("div",{className:"advanced-row"},r.a.createElement(b.a,{variant:"contained",color:"primary",onClick:this.handleSubmit,style:{width:"100%"}},"Save")))))}}]),n}(r.a.Component);c.a.render(r.a.createElement(r.a.StrictMode,null,r.a.createElement(B,null)),document.getElementById("root"))},73:function(e,t,n){e.exports=n(127)},79:function(e,t,n){}},[[73,1,2]]]);
//# sourceMappingURL=main.23a6cae1.chunk.js.map