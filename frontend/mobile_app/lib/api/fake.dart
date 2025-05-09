import 'package:mobile_app/api/api.dart';
import 'package:mobile_app/api/common.dart';
import 'package:mobile_app/api/exceptions.dart';
import 'package:mobile_app/api/requests.dart';
import 'package:mobile_app/api/responses.dart';

const Duration delay = Duration(milliseconds: 500);

class FakeServer implements Api {
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
  ) {
    // TODO: implement getCompetitionsCompetitors
    throw UnimplementedError();
  }

  @override
  Future<List<IndividualGroup>> getCompetitionsIndividualGroups(
    int competitionId,
  ) {
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
  Future<List<CompetitorGroupDetail>> getIndividualGroupCompetitors(
    int groupId,
  ) {
    // TODO: implement getIndividualGroupCompetitors
    throw UnimplementedError();
  }

  @override
  Future<FinalGrid> getIndividualGroupFinalGrid(int groupId) {
    defineSparring(sparringId, top, bottom) => Sparring(
      sparringId,
      SparringPlace(
        top.id,
        top.shrink(),
        RangeGroup(1, 3, 3, RangeType.one2ten, [], 0),
        false,
        null,
        0,
      ),
      bottom != null
          ? SparringPlace(
            bottom.id,
            bottom.shrink(),
            RangeGroup(1, 3, 3, RangeType.one2ten, [], 0),
            false,
            null,
            0,
          )
          : null,
      SparringState.ongoing,
    );
    List<Sparring> sparrings = [
      defineSparring(1, lebedev, piyavkin),
      defineSparring(2, kozakova, dudkina),
      defineSparring(3, kravchenko, null),
      defineSparring(4, demidenko, novokhatskiy),

      defineSparring(5, lebedev, kozakova),
      defineSparring(6, novokhatskiy, kravchenko),

      defineSparring(7, kozakova, kravchenko),

      defineSparring(8, lebedev, novokhatskiy),
    ];
    return Future.delayed(
      delay,
      () => switch (groupId) {
        1 => FinalGrid(
          1,
          Quarterfinal(sparrings[0], sparrings[1], sparrings[2], sparrings[3]),
          Semifinal(sparrings[4], sparrings[5]),
          Final(sparrings[6], sparrings[7]),
        ),
        _ => throw NotFoundException("Группа или сетка не найдена"),
      },
    );
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
  Future<QualificationRoundFull> getQualificationSectionsRound(
    int sectionId,
    int roundOrdinal,
  ) {
    // TODO: implement getQualificationSectionsRound
    throw UnimplementedError();
  }

  @override
  Future<RangeGroup> getQualificationSectionsRoundsRanges(
    int sectionId,
    int roundOrdinal,
  ) {
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

  static final CompetitorFull lebedev = CompetitorFull(
    1,
    "Лебедев Антон",
    DateTime(2000),
    Gender.male,
    BowClass.classic,
    SportsRank.masterInternational,
    "Минск",
    null,
    null,
  );
  static final CompetitorFull piyavkin = CompetitorFull(
    2,
    "Пиявкин Антон",
    DateTime(2000),
    Gender.male,
    BowClass.block,
    SportsRank.master,
    null,
    "Федерация водоплавающих",
    null,
  );
  static final CompetitorFull kozakova = CompetitorFull(
    3,
    "Козакова Анна",
    DateTime(2000),
    Gender.female,
    BowClass.classic3D,
    SportsRank.candidateForMaster,
    null,
    null,
    "Клуб go",
  );
  static final CompetitorFull dudkina = CompetitorFull(
    4,
    "Дудкина София",
    DateTime(2000),
    Gender.female,
    BowClass.classicNewbie,
    SportsRank.firstClass,
    null,
    "Федерация бекенда",
    "Клуб go",
  );
  static final CompetitorFull kravchenko = CompetitorFull(
    5,
    "Кравченко Никита",
    DateTime(2000),
    Gender.male,
    BowClass.compound3D,
    SportsRank.meritedMaster,
    "Россия",
    "Федерация фулстека",
    "Клуб js",
  );
  static final CompetitorFull demidenko = CompetitorFull(
    6,
    "Демиденко Никита",
    DateTime(2000),
    Gender.male,
    BowClass.long3D,
    SportsRank.secondClass,
    "СПБ",
    null,
    "Клуб rust",
  );
  static final CompetitorFull novokhatskiy = CompetitorFull(
    7,
    "Новохацкий Данил",
    DateTime(2000),
    Gender.male,
    BowClass.peripheral,
    SportsRank.thirdClass,
    "Владивосток",
    "Федерация фронтеда",
    null,
  );
}
