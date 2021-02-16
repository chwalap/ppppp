package ppppp

import scala.concurrent.duration._

import io.gatling.core.Predef._
import io.gatling.http.Predef._
import io.gatling.jdbc.Predef._

class ppppp extends Simulation {

	val httpProtocol = http
		.baseUrl("http://localhost:8080")
		.inferHtmlResources()
		.acceptHeader("text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
		.acceptEncodingHeader("gzip, deflate")
		.acceptLanguageHeader("pl,en-US;q=0.9,en;q=0.8")
		.doNotTrackHeader("1")
		.userAgentHeader("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.150 Safari/537.36")

	val headers_0 = Map(
		"Cache-Control" -> "max-age=0",
		"Proxy-Connection" -> "keep-alive",
		"Sec-Fetch-Dest" -> "document",
		"Sec-Fetch-Mode" -> "navigate",
		"Sec-Fetch-Site" -> "none",
		"Sec-Fetch-User" -> "?1",
		"Upgrade-Insecure-Requests" -> "1",
		"sec-ch-ua" -> """"Chromium";v="88", "Google Chrome";v="88", ";Not A Brand";v="99"""",
		"sec-ch-ua-mobile" -> "?0")

	val headers_1 = Map(
		"Accept" -> "image/avif,image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8",
		"Cache-Control" -> "no-cache",
		"Pragma" -> "no-cache",
		"Proxy-Connection" -> "keep-alive",
		"Sec-Fetch-Dest" -> "image",
		"Sec-Fetch-Mode" -> "no-cors",
		"Sec-Fetch-Site" -> "same-origin",
		"sec-ch-ua" -> """"Chromium";v="88", "Google Chrome";v="88", ";Not A Brand";v="99"""",
		"sec-ch-ua-mobile" -> "?0")

	val headers_2 = Map(
		"Cache-Control" -> "max-age=0",
		"Origin" -> "http://localhost:8080",
		"Proxy-Connection" -> "keep-alive",
		"Sec-Fetch-Dest" -> "document",
		"Sec-Fetch-Mode" -> "navigate",
		"Sec-Fetch-Site" -> "same-origin",
		"Sec-Fetch-User" -> "?1",
		"Upgrade-Insecure-Requests" -> "1",
		"sec-ch-ua" -> """"Chromium";v="88", "Google Chrome";v="88", ";Not A Brand";v="99"""",
		"sec-ch-ua-mobile" -> "?0")

	val headers_7 = Map(
		"Proxy-Connection" -> "keep-alive",
		"Sec-Fetch-Dest" -> "document",
		"Sec-Fetch-Mode" -> "navigate",
		"Sec-Fetch-Site" -> "same-origin",
		"Sec-Fetch-User" -> "?1",
		"Upgrade-Insecure-Requests" -> "1",
		"sec-ch-ua" -> """"Chromium";v="88", "Google Chrome";v="88", ";Not A Brand";v="99"""",
		"sec-ch-ua-mobile" -> "?0")

	val headers_10 = Map(
		"Cache-Control" -> "max-age=0",
		"Proxy-Connection" -> "keep-alive",
		"Sec-Fetch-Dest" -> "document",
		"Sec-Fetch-Mode" -> "navigate",
		"Sec-Fetch-Site" -> "same-origin",
		"Sec-Fetch-User" -> "?1",
		"Upgrade-Insecure-Requests" -> "1",
		"sec-ch-ua" -> """"Chromium";v="88", "Google Chrome";v="88", ";Not A Brand";v="99"""",
		"sec-ch-ua-mobile" -> "?0")



	val scn = scenario("ppppp")
		.exec(http("ppppp_0:GET_http://localhost:8080/signup")
			.get("/signup")
			.headers(headers_0)
			.resources(http("ppppp_1:GET_http://localhost:8080/resources/favicon.ico")
			.get("/resources/favicon.ico")
			.headers(headers_1)
			.check(bodyBytes.is(RawFileBody("ppppp/ppppp/0001_response.dat"))))
			.check(bodyBytes.is(RawFileBody("ppppp/ppppp/0000_response.html"))))
		.pause(7)
		.exec(http("ppppp_2:POST_http://localhost:8080/signup")
			.post("/signup")
			.headers(headers_2)
			.formParam("username", "dupa")
			.formParam("password", "dupa"))
		.pause(3)
		.exec(http("ppppp_3:POST_http://localhost:8080/login")
			.post("/login")
			.headers(headers_2)
			.formParam("username", "dupa")
			.formParam("password", "dupa"))
		.pause(13)
		.exec(http("ppppp_4:POST_http://localhost:8080/worker/new")
			.post("/worker/new")
			.headers(headers_2)
			.formParam("city", "3090759")
			.formParam("interval", "10"))
		.pause(6)
		.exec(http("ppppp_5:POST_http://localhost:8080/worker/new")
			.post("/worker/new")
			.headers(headers_2)
			.formParam("city", "993800")
			.formParam("interval", "15"))
		.pause(6)
		.exec(http("ppppp_6:POST_http://localhost:8080/worker/new")
			.post("/worker/new")
			.headers(headers_2)
			.formParam("city", "4140963")
			.formParam("interval", "5"))
		.pause(7)
		.exec(http("ppppp_7:GET_http://localhost:8080/worker/edit?id=182")
			.get("/worker/edit?id=182")
			.headers(headers_7)
			.check(bodyBytes.is(RawFileBody("ppppp/ppppp/0007_response.html"))))
		.pause(5)
		.exec(http("ppppp_8:POST_http://localhost:8080/worker/edit")
			.post("/worker/edit")
			.headers(headers_2)
			.formParam("id", "182")
			.formParam("running", "1")
			.formParam("city", "2147714")
			.formParam("interval", "3"))
		.pause(5)
		.exec(http("ppppp_9:GET_http://localhost:8080/worker/chart?id=182")
			.get("/worker/chart?id=182")
			.headers(headers_7)
			.check(bodyBytes.is(RawFileBody("ppppp/ppppp/0009_response.html"))))
		.pause(5)
		.exec(http("ppppp_10:GET_http://localhost:8080/worker/chart?id=182")
			.get("/worker/chart?id=182")
			.headers(headers_10)
			.resources(http("ppppp_11:GET_http://localhost:8080/resources/favicon.ico")
			.get("/resources/favicon.ico")
			.headers(headers_1)
			.check(bodyBytes.is(RawFileBody("ppppp/ppppp/0011_response.dat"))))
			.check(bodyBytes.is(RawFileBody("ppppp/ppppp/0010_response.html"))))
		.pause(18)
		.exec(http("ppppp_12:GET_http://localhost:8080/worker/chart?id=182")
			.get("/worker/chart?id=182")
			.headers(headers_7)
			.check(bodyBytes.is(RawFileBody("ppppp/ppppp/0012_response.html"))))
		.pause(25)
		.exec(http("ppppp_13:GET_http://localhost:8080/worker/delete?id=182")
			.get("/worker/delete?id=182")
			.headers(headers_7)
			.check(bodyBytes.is(RawFileBody("ppppp/ppppp/0013_response.html"))))
		.pause(5)
		.exec(http("ppppp_14:GET_http://localhost:8080/worker/pause?id=184")
			.get("/worker/pause?id=184")
			.headers(headers_7)
			.check(bodyBytes.is(RawFileBody("ppppp/ppppp/0014_response.html"))))
		.pause(4)
		.exec(http("ppppp_15:GET_http://localhost:8080/worker/pause?id=183")
			.get("/worker/pause?id=183")
			.headers(headers_7)
			.check(bodyBytes.is(RawFileBody("ppppp/ppppp/0015_response.html"))))
		.pause(1)
		.exec(http("ppppp_16:GET_http://localhost:8080/worker/start?id=183")
			.get("/worker/start?id=183")
			.headers(headers_7)
			.check(bodyBytes.is(RawFileBody("ppppp/ppppp/0016_response.html"))))
		.pause(8)
		.exec(http("ppppp_17:POST_http://localhost:8080/logout")
			.post("/logout")
			.headers(headers_2))

	setUp(scn.inject(atOnceUsers(1))).protocols(httpProtocol)
}