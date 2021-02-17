package ppppp

import scala.concurrent.duration._

import io.gatling.core.Predef._
import io.gatling.http.Predef._
import io.gatling.jdbc.Predef._

class FiveP extends Simulation {

	val httpProtocol = http
		.baseUrl("http://host.docker.internal:8080")
		.inferHtmlResources()
		.acceptHeader("text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
		.acceptEncodingHeader("gzip, deflate")
		.acceptLanguageHeader("pl,en-US;q=0.9,en;q=0.8")
		.doNotTrackHeader("1")
		.userAgentHeader("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.150 Safari/537.36")

	val headers_0 = Map(
		"Cache-Control" -> "max-age=0",
		"Proxy-Connection" -> "keep-alive",
		"Upgrade-Insecure-Requests" -> "1")

	val headers_1 = Map(
		"Accept" -> "text/css,*/*;q=0.1",
		"If-Modified-Since" -> "Wed, 17 Feb 2021 17:07:57 GMT",
		"Proxy-Connection" -> "keep-alive")

	val headers_2 = Map(
		"Accept" -> "image/avif,image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8",
		"Cache-Control" -> "no-cache",
		"Pragma" -> "no-cache",
		"Proxy-Connection" -> "keep-alive")

	val headers_3 = Map(
		"Cache-Control" -> "max-age=0",
		"Origin" -> "http://host.docker.internal:8080",
		"Proxy-Connection" -> "keep-alive",
		"Upgrade-Insecure-Requests" -> "1")

	val headers_7 = Map(
		"Proxy-Connection" -> "keep-alive",
		"Upgrade-Insecure-Requests" -> "1")

		object Login {

			val login = exec(http("Login Page")
				.get("/login")
				.headers(headers_0)
				.resources(http("GET styles.css")
				.get("/resources/styles.css")
				.headers(headers_1),
				http("GET favicon.ico")
				.get("/resources/favicon.ico")
				.headers(headers_2)
				.check(bodyBytes.is(RawFileBody("5p/5p/0002_response.ico"))))
				.check(bodyBytes.is(RawFileBody("5p/5p/0000_response.html"))))
			.pause(4)
			.exec(http("Login Page")
				.post("/login")
				.headers(headers_3)
				.formParam("username", "q")
				.formParam("password", "q"))
			.pause(6)
		}

		object NewWorker {

			val newWorker = exec(http("New worker")
			.post("/worker/new")
			.headers(headers_3)
			.formParam("city", "538560")
			.formParam("interval", "1"))
		.pause(2)
		.exec(http("New worker")
			.post("/worker/new")
			.headers(headers_3)
			.formParam("city", "2267057")
			.formParam("interval", "60")
			.resources(http("New worker")
			.post("/worker/new")
			.headers(headers_3)
			.formParam("city", "2643743")
			.formParam("interval", "60")))
		.pause(5)

		}


		object DeleteWorker {

			val deleteWorker = exec(http("Delete Worker")
				.get("/worker/delete?id=41")
				.headers(headers_7)
				.check(bodyBytes.is(RawFileBody("5p/5p/0007_response.html"))))
			.pause(1)
			.exec(http("Delete Worker")
				.get("/worker/delete?id=40")
				.headers(headers_7)
				.check(bodyBytes.is(RawFileBody("5p/5p/0008_response.html"))))
			.pause(1)
			.exec(http("Delete Worker")
				.get("/worker/delete?id=39")
				.headers(headers_7)
				.check(bodyBytes.is(RawFileBody("5p/5p/0009_response.html"))))
		}

		object EditWorker {

			val editWorker = exec(http("Edit Page")
			.get("/worker/edit?id=9")
			.headers(headers_7)
			.resources(http("GET styles.css")
			.get("/resources/styles.css")
			.headers(headers_1))
			.check(bodyBytes.is(RawFileBody("5p/5p/0014_response.html"))))
		.pause(3)
		.exec(http("Edit Page")
			.post("/worker/edit")
			.headers(headers_3)
			.formParam("id", "9")
			.formParam("running", "1")
			.formParam("city", "2643743")
			.formParam("interval", "81"))
		.pause(2)

		}

		object PauseWorker {

			val pauseWorker = exec(http("Pause Worker")
			.get("/worker/pause?id=37")
			.headers(headers_7)
			.check(bodyBytes.is(RawFileBody("5p/5p/0010_response.html"))))
		.pause(2)
		.exec(http("Pause Worker")
			.get("/worker/pause?id=22")
			.headers(headers_7)
			.check(bodyBytes.is(RawFileBody("5p/5p/0011_response.html"))))
		.pause(1)
		.exec(http("Pause Worker")
			.get("/worker/pause?id=5")
			.headers(headers_7)
			.check(bodyBytes.is(RawFileBody("5p/5p/0012_response.html"))))
		.pause(1)
		.exec(http("Pause Worker")
			.get("/worker/pause?id=15")
			.headers(headers_7)
			.check(bodyBytes.is(RawFileBody("5p/5p/0013_response.html"))))
			.pause(2)
		}

		object ViewChart {
		
		val viewChart = exec(http("View Chart")
			.get("/worker/chart?id=5")
			.headers(headers_7)
			.check(bodyBytes.is(RawFileBody("5p/5p/0017_response.html"))))
		.pause(4)
		.exec(http("View Chart")
			.get("/worker/chart?id=5&date=2021-02-17")
			.headers(headers_7)
			.check(bodyBytes.is(RawFileBody("5p/5p/0018_response.html"))))
		.pause(17)
		}


		object Logout {

			val logout = exec(http("Logout")
			.get("/logout")
			.headers(headers_7)
			.check(bodyBytes.is(RawFileBody("5p/5p/0019_response.html"))))
		}

		val Basic = scenario("Basic").exec(Login.login, NewWorker.newWorker, ViewChart.viewChart, Logout.logout)
		val Advanced = scenario("Advanced").exec(Login.login, NewWorker.newWorker, EditWorker.editWorker, ViewChart.viewChart, Logout.logout)
		val Pause = scenario("Just Pausing...").exec(Login.login, PauseWorker.pauseWorker, Logout.logout)


		setUp(
			Basic.inject(rampUsers(5000).during(300.seconds)),
			Advanced.inject(rampUsers(5000).during(300.seconds)),
			Pause.inject(atOnceUsers(1))
		).protocols(httpProtocol)
}