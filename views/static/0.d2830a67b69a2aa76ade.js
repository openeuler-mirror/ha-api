webpackJsonp([0],{1384:function(e,t,n){"use strict";function r(e){return e&&e.__esModule?e:{default:e}}Object.defineProperty(t,"__esModule",{value:!0}),t.getLicense=t.login=void 0;var u=n(119),a=r(u),s=n(286),c=r(s),i=(t.login=function(){var e=(0,c.default)(a.default.mark(function e(t){return a.default.wrap(function(e){for(;;)switch(e.prev=e.next){case 0:return e.abrupt("return",(0,i.request)({url:l,method:"post",data:t}));case 1:case"end":return e.stop()}},e,this)}));return function(t){return e.apply(this,arguments)}}(),t.getLicense=function(){var e=(0,c.default)(a.default.mark(function e(t){return a.default.wrap(function(e){for(;;)switch(e.prev=e.next){case 0:return e.abrupt("return",(0,i.request)({url:f,method:"get",data:t}));case 1:case"end":return e.stop()}},e,this)}));return function(t){return e.apply(this,arguments)}}(),n(27)),o=i.config.api,l=o.userLogin,f=o.license},613:function(e,t,n){"use strict";function r(e){return e&&e.__esModule?e:{default:e}}Object.defineProperty(t,"__esModule",{value:!0});var u=n(3),a=r(u),s=n(190),c=r(s),i=n(119),o=r(i),l=n(288),f=n(1384),d=n(27),p=n(102);t.default={namespace:"login",state:{license:JSON.parse(window.localStorage.getItem(p.prefix+"license"))},effects:{login:o.default.mark(function e(t,n){var r=t.payload,u=n.put,a=n.call;return o.default.wrap(function(e){for(;;)switch(e.prev=e.next){case 0:return e.next=2,a(f.login,r);case 2:return window.localStorage.setItem(p.prefix+"username",r.username),e.next=5,u(l.routerRedux.push("/resource"));case 5:case"end":return e.stop()}},e,this)}),getLicense:o.default.mark(function e(t,n){var r,u,a=n.call,s=n.put;return o.default.wrap(function(e){for(;;)switch(e.prev=e.next){case 0:return e.next=2,a(f.getLicense);case 2:return r=e.sent,u=(0,d.getLicenseType)(r.data),e.next=6,s({type:"updateLicense",payload:{license:u}});case 6:case"end":return e.stop()}},e,this)})},reducers:{updateLicense:function(e,t){var n=t.payload;return window.localStorage.setItem(p.prefix+"license",(0,c.default)(n.license)),(0,a.default)({},e,n)}}},e.exports=t.default}});