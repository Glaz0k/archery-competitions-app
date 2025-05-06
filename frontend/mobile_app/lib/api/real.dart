import 'package:mobile_app/api/api.dart';
import 'package:mobile_app/api/requests.dart';
import 'package:mobile_app/api/responses.dart';

class RealServer implements Api {
  @override
  Future<CompetitorCompetitionDetail> changeCompetitorStatus(int competitionId, int competitorId, bool status) {
    // TODO: implement changeCompetitorStatus
    throw UnimplementedError();
  }

  @override
  Future<Range> endQualificationSectionsRoundsRange(int sectionId, int roundOrdinal, int rangeOrdinal) {
    // TODO: implement endQualificationSectionsRoundsRange
    throw UnimplementedError();
  }

  @override
  Future<Range> endSparringPlacesRange(int placeId, int rangeOrdinal) {
    // TODO: implement endSparringPlacesRange
    throw UnimplementedError();
  }

  @override
  Future<List<CompetitorCompetitionDetail>> getCompetitionsCompetitors(int competitionId) {
    // TODO: implement getCompetitionsCompetitors
    throw UnimplementedError();
  }

  @override
  Future<List<IndividualGroup>> getCompetitionsIndividualGroups(int competitionId) {
    // TODO: implement getCompetitionsIndividualGroups
    throw UnimplementedError();
  }

  @override
  Future<CompetitorFull> getCompetitor(int competitorId) {
    // TODO: implement getCompetitor
    throw UnimplementedError();
  }

  @override
  Future<Cup> getCup(int cupId) {
    // TODO: implement getCup
    throw UnimplementedError();
  }

  @override
  Future<List<Cup>> getCups() {
    // TODO: implement getCups
    throw UnimplementedError();
  }

  @override
  Future<List<Competition>> getCupsCompetitions(int cupId) {
    // TODO: implement getCupsCompetitions
    throw UnimplementedError();
  }

  @override
  Future<IndividualGroup> getIndividualGroup(int groupId) {
    // TODO: implement getIndividualGroup
    throw UnimplementedError();
  }

  @override
  Future<List<CompetitorGroupDetail>> getIndividualGroupCompetitors(int groupId) {
    // TODO: implement getIndividualGroupCompetitors
    throw UnimplementedError();
  }

  @override
  Future<FinalGrid> getIndividualGroupFinalGrid(int groupId) {
    // TODO: implement getIndividualGroupFinalGrid
    throw UnimplementedError();
  }

  @override
  Future<QualificationTable> getIndividualGroupQualificationTable(int groupId) {
    // TODO: implement getIndividualGroupQualificationTable
    throw UnimplementedError();
  }

  @override
  Future<Section> getQualificationSection(int sectionId) {
    // TODO: implement getQualificationSection
    throw UnimplementedError();
  }

  @override
  Future<QualificationRoundFull> getQualificationSectionsRound(int sectionId, int roundOrdinal) {
    // TODO: implement getQualificationSectionsRound
    throw UnimplementedError();
  }

  @override
  Future<RangeGroup> getQualificationSectionsRoundsRanges(int sectionId, int roundOrdinal) {
    // TODO: implement getQualificationSectionsRoundsRanges
    throw UnimplementedError();
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
  Future<void> login(Credentials credentials) {
    // TODO: implement login
    throw UnimplementedError();
  }

  @override
  Future<void> logout() {
    // TODO: implement logout
    throw UnimplementedError();
  }

  @override
  Future<CompetitorFull> putCompetitor(int competitorId, ChangeCompetitor request) {
    // TODO: implement putCompetitor
    throw UnimplementedError();
  }

  @override
  Future<Range> putQualificationSectionsRoundsRange(int sectionId, int roundOrdinal, ChangeRange request) {
    // TODO: implement putQualificationSectionsRoundsRange
    throw UnimplementedError();
  }

  @override
  Future<Range> putSparringPlacesRange(int placeId, ChangeRange request) {
    // TODO: implement putSparringPlacesRange
    throw UnimplementedError();
  }

  @override
  Future<ShootOut> putSparringPlacesShootOut(int placeId, ChangeShootOut request) {
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