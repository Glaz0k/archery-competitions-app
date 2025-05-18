import 'dart:convert';

import 'package:http/http.dart' as http;
import 'package:mobile_app/api/api.dart';
import 'package:mobile_app/api/exceptions.dart';
import 'package:mobile_app/api/requests.dart';
import 'package:mobile_app/api/responses.dart';

const String backend = "example.com";

class RealServer implements Api {
  final http.Client client = http.Client();

  void validate(
    http.Response response, {
    String? notFoundMessage,
    String? invalidParametersMessage,
    String? alreadyExistsMessage,
    String? badActionMessage,
    String? invalidScoreMessage,
  }) {
    var defaultMessage = "Получена некорректная ошибка от сервера";
    switch (response.statusCode) {
      case 404:
        throw NotFoundException(notFoundMessage!);

      case 400:
        Map<String, dynamic> body = jsonDecode(response.body);
        switch (body["error"]) {
          case "INVALID PARAMETERS":
            throw InvalidParametersException(
              invalidParametersMessage ?? defaultMessage,
            );
          case "EXISTS":
            throw AlreadyExistException(alreadyExistsMessage ?? defaultMessage);
          case "BAD ACTION":
            throw BadActionException(badActionMessage ?? defaultMessage);
          case "INVALID SCORE":
            Map<String, dynamic> details = body["details"];
            throw InvalidScoreException(
              invalidScoreMessage ?? defaultMessage,
              details["shot_ordinal"],
              details["type"],
            );
        }
      case 200 || 204:
        return;
    }
  }

  @override
  Future<CompetitorCompetitionDetail> changeCompetitorStatus(
    int competitionId,
    int competitorId,
    bool status,
  ) {
    // TODO: implement changeCompetitorStatus
    throw UnimplementedError();
  }

  @override
  Future<Range> endQualificationSectionsRoundsRange(
    int sectionId,
    int roundOrdinal,
    int rangeOrdinal,
  ) {
    // TODO: implement endQualificationSectionsRoundsRange
    throw UnimplementedError();
  }

  @override
  Future<Range> endSparringPlacesRange(int placeId, int rangeOrdinal) {
    // TODO: implement endSparringPlacesRange
    throw UnimplementedError();
  }

  @override
  Future<List<CompetitorCompetitionDetail>> getCompetitionsCompetitors(
    int competitionId,
  ) async {
    final response = await client.get(
      Uri.https(backend, "/competitions/$competitionId/competitors"),
    );
    validate(response, notFoundMessage: "Соревнование не найдено");
    final List<dynamic> jsonList = jsonDecode(response.body);
    return jsonList
        .map((json) => CompetitorCompetitionDetail.fromJson(json))
        .toList();
  }

  @override
  Future<List<IndividualGroup>> getCompetitionsIndividualGroups(
    int competitionId,
  ) async {
    final response = await client.get(
      Uri.https(backend, "/competitions/$competitionId/individual_groups"),
    );
    validate(response, notFoundMessage: "Соревнование не найдено");
    final List<dynamic> jsonList = jsonDecode(response.body);
    return jsonList.map((json) => IndividualGroup.fromJson(json)).toList();
  }

  @override
  Future<CompetitorFull> getCompetitor(int competitorId) async {
    final response = await client.get(
      Uri.https(backend, "/competitors/$competitorId"),
    );
    validate(response, notFoundMessage: "Участник не найден");
    return CompetitorFull.fromJson(jsonDecode(response.body));
  }

  @override
  Future<Cup> getCup(int cupId) async {
    final response = await client.get(Uri.https(backend, "/cups/$cupId"));
    validate(response, notFoundMessage: "Кубок не найден");
    return Cup.fromJson(jsonDecode(response.body));
  }

  @override
  Future<List<Cup>> getCups() async {
    final response = await client.get(Uri.https(backend, "/cups"));
    validate(response, notFoundMessage: "Кубки не найдены");
    final List<dynamic> jsonList = jsonDecode(response.body);
    return jsonList.map((json) => Cup.fromJson(json)).toList();
  }

  @override
  Future<List<Competition>> getCupsCompetitions(int cupId) async {
    final response = await client.get(
      Uri.https(backend, "/cups/$cupId/competitions"),
    );
    validate(response, notFoundMessage: "Кубок не найден");
    final List<dynamic> jsonList = jsonDecode(response.body);
    return jsonList.map((json) => Competition.fromJson(json)).toList();
  }

  @override
  Future<IndividualGroup> getIndividualGroup(int groupId) async {
    final response = await client.get(
      Uri.https(backend, "/individual_groups/$groupId"),
    );
    validate(response, notFoundMessage: "Группа не найдена");
    return IndividualGroup.fromJson(jsonDecode(response.body));
  }

  @override
  Future<List<CompetitorGroupDetail>> getIndividualGroupCompetitors(
    int groupId,
  ) async {
    final response = await client.get(
      Uri.https(backend, "/individual_groups/$groupId/competitors"),
    );
    validate(response, notFoundMessage: "Группа не найдена");
    final List<dynamic> jsonList = jsonDecode(response.body);
    return jsonList
        .map((json) => CompetitorGroupDetail.fromJson(json))
        .toList();
  }

  @override
  Future<FinalGrid> getIndividualGroupFinalGrid(int groupId) async {
    var response = await client.get(
      Uri.https(backend, "/individual_groups/$groupId/final_grid"),
    );
    validate(response, notFoundMessage: "Группа или сетка не найдена");
    return FinalGrid.fromJson(jsonDecode(response.body));
  }

  @override
  Future<QualificationTable> getIndividualGroupQualificationTable(
    int groupId,
  ) async {
    var response = await client.get(
      Uri.https(backend, "/individual_groups/$groupId/qualification"),
    );
    validate(response);
    return QualificationTable.fromJson(jsonDecode(response.body));
  }

  @override
  Future<Section> getQualificationSection(int sectionId) async {
    var response = await client.get(
      Uri.https(backend, "/qualification_sections/$sectionId"),
    );
    validate(response, notFoundMessage: "Секция не найдена");
    return Section.fromJson(jsonDecode(response.body));
  }

  @override
  Future<QualificationRoundFull> getQualificationSectionsRound(
    int sectionId,
    int roundOrdinal,
  ) async {
    var response = await client.get(
      Uri.https(
        backend,
        "/qualification_sections/$sectionId/rounds/$roundOrdinal",
      ),
    );
    validate(response);
    return QualificationRoundFull.fromJson(jsonDecode(response.body));
  }

  @override
  Future<RangeGroup> getQualificationSectionsRoundsRanges(
    int sectionId,
    int roundOrdinal,
  ) async {
    var response = await client.get(
      Uri.https(
        backend,
        "/qualification_sections/$sectionId/rounds/$roundOrdinal/ranges",
      ),
    );
    validate(response);
    return RangeGroup.fromJson(jsonDecode(response.body));
  }

  @override
  Future<SparingPlace> getSparringPlace(int placeId) {
    // TODO: implement getSparringPlace
    throw UnimplementedError();
  }

  @override
  Future<RangeGroup> getSparringPlacesRanges(int placeId) {
    // TODO: implement getSparringPlacesRanges
    throw UnimplementedError();
  }

  @override
  Future<int> login(Credentials credentials) async {
    // TODO: эту фигню ещё не заимплементить
    throw UnimplementedError();
  }

  @override
  Future<void> logout() {
    // TODO: implement logout
    throw UnimplementedError();
  }

  @override
  Future<CompetitorFull> putCompetitor(
    int competitorId,
    ChangeCompetitor request,
  ) {
    // TODO: implement putCompetitor
    throw UnimplementedError();
  }

  @override
  Future<Range> putQualificationSectionsRoundsRange(
    int sectionId,
    int roundOrdinal,
    ChangeRange request,
  ) {
    // TODO: implement putQualificationSectionsRoundsRange
    throw UnimplementedError();
  }

  @override
  Future<Range> putSparringPlacesRange(int placeId, ChangeRange request) {
    // TODO: implement putSparringPlacesRange
    throw UnimplementedError();
  }

  @override
  Future<ShootOut> putSparringPlacesShootOut(
    int placeId,
    ChangeShootOut request,
  ) {
    // TODO: implement putSparringPlacesShootOut
    throw UnimplementedError();
  }

  @override
  Future<void> register(Credentials credentials) {
    // TODO: implement register
    throw UnimplementedError();
  }

  @override
  Future<CompetitorFull> registerCompetitor(ChangeCompetitor request) {
    // TODO: implement registerCompetitor
    throw UnimplementedError();
  }
}
