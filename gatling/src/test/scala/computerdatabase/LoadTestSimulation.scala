package computerdatabase

import io.gatling.core.Predef._
import io.gatling.http.Predef._
import io.gatling.core.structure.PopulationBuilder

import scala.concurrent.duration._
import scala.util.parsing.json._
import scala.util.Random

import java.security.SecureRandom

class LoadTestSimulation extends Simulation {
  private val serverURI = "http://localhost:28080/v1"

  object User {
    def create() = exec(
      http("create_user")
        .post(f"$serverURI/users")
        .check(
          status.is(201),
          jsonPath("$.id").find.saveAs("user_id"),
          jsonPath("$.name").find.saveAs("user_name")
        )
    ).exitHereIfFailed

    def refer(userID: String, userName: String) = exec(
      http("refer_user")
        .get(f"$serverURI/users/$userID")
        .check(
          status.is(200),
          jsonPath("$.id").find.is(userID),
          jsonPath("$.name").find.is(userName),
          jsonPath("$.secret_key").find.saveAs("secret_key")
        ),
    ).exitHereIfFailed
  }

  object Page {
    private val param = (userID: String, name: String, secretKey: String) => {
      s"""{
        "user_id": "$userID",
        "name": "$name",
        "secret_key": "$secretKey"
      }"""
    }

    private def redirectRequest(requestTitle: String, param: String) = http(requestTitle)
      .post(f"$serverURI/display")
      .body(StringBody(param))
      .asJson

    def redirectStatus(userID: String, name: String, secretKey: String) = exec(
      redirectRequest("redirect_status", param(userID, name, secretKey))
        .disableFollowRedirect
        .check(status.is(301))
    ).exitHereIfFailed

    def redirectPage(userID: String, name: String, secretKey: String) = exec(
      redirectRequest("redirect_page", param(userID, name, secretKey))
        .check(status.is(200))
    ).exitHereIfFailed
  }

  def execScenario(name: String) = {
    scenario(name)
      .exec(User.create())
      .exec(User.refer("${user_id}", "${user_name}"))
      .exec(Page.redirectStatus("${user_id}", "${user_name}", "${secret_key}"))
      .exec(Page.redirectPage("${user_id}", "${user_name}", "${secret_key}"))
  }

  def getURI(): String = {
    return serverURI
  }

  setUp(
    execScenario("Load test base simulation").inject(
      atOnceUsers(1),
    ).protocols(http.baseUrl(getURI())),
  )
}


class LoadTestSimulationRampUp10Users extends LoadTestSimulation {
  override def setUp(populationBuilders: PopulationBuilder*): SetUp =
    super.setUp(
      execScenario("Load test simulation: rampUp 10 users")
        .inject(rampUsers(10) during (5 seconds))
        .protocols(http.baseUrl(getURI()))
    )
}

class LoadTestSimulationCompoundTest extends LoadTestSimulation {
  override def setUp(populationBuilders: PopulationBuilder*): SetUp =
    super.setUp(
      execScenario("Load test simulation: init 10 users, rampUp 10 users(5 sec), constant 30 users(5 sec)")
        .inject(
          atOnceUsers(10),
          rampUsers(10) during (5 seconds),
          rampUsers(30) during (5 seconds)
        )
        .protocols(http.baseUrl(getURI()))
    )
}
