import 'package:mobile_app/api/api.dart';
import 'package:mobile_app/api/common.dart';
import 'package:mobile_app/api/exceptions.dart';
import 'package:mobile_app/api/requests.dart';
import 'package:mobile_app/api/responses.dart';

const Duration delay = Duration(milliseconds: 500);

class FakeServer implements Api {
  final Map<int, Cup> _cups = {
    1: Cup(1, "World Archery Championship", "Berlin, Germany", "2023-2024"),
    2: Cup(2, "European Archery Cup", "Paris, France", "2023"),
  };

  final Section _section = Section(
    1,
    novokhatskiy.shrink(),
    1,
    [],
    100,
    23,
    24,
    SportsRank.masterInternational,
  );

  final Map<int, Competition> _competitions = {
    1: Competition(
      1,
      CompetitionStage.I,
      DateTime(2025, 2, 3),
      DateTime(2025, 2, 17),
      false,
    ),
    2: Competition(
      2,
      CompetitionStage.II,
      DateTime(2025, 3, 22),
      DateTime(2025, 3, 28),
      true,
    ),
  };

  final RangeGroup _rangeGroup = RangeGroup(1, 3, 3, RangeType.one2ten, [
    Range(1, 1, false, [Shot(1, "4"), Shot(2, "3"), Shot(3, "M")], 7),
    Range(1, 1, true, [Shot(1, "X"), Shot(2, null), Shot(3, null)], null),
  ], 17);

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
    var range = _rangeGroup.ranges[rangeOrdinal - 1];
    range.isActive = false;
    return Future.delayed(delay, () => range);
  }

  @override
  Future<Range> endSparringPlacesRange(int placeId, int rangeOrdinal) {
    var range = _rangeGroup.ranges[rangeOrdinal - 1];
    range.isActive = false;
    return Future.delayed(delay, () => range);
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
  Future<Cup> getCup(int cupId) async {
    if (_cups.containsKey(cupId)) {
      return _cups[cupId]!;
    } else {
      throw NotFoundException("Кубок не найден");
    }
  }

  @override
  Future<List<Cup>> getCups() {
    if (_cups.isEmpty) {
      throw NotFoundException("Кубки не найдены");
    }
    List<Cup> cups = _cups.values.toList();
    return Future.delayed(delay, () => cups);
  }

  @override
  Future<List<Competition>> getCupsCompetitions(int cupId) {
    if (_competitions.isEmpty) {
      throw NotFoundException("Соревнования не найдены");
    }
    List<Competition> competitions = _competitions.values.toList();

    return Future.delayed(delay, () => competitions);
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
    return Future.delayed(delay, () {
      switch (groupId) {
        case 1:
          return [
            CompetitorGroupDetail(1, lebedev),
            CompetitorGroupDetail(1, piyavkin),
            CompetitorGroupDetail(1, kozakova),
            CompetitorGroupDetail(1, dudkina),
            CompetitorGroupDetail(1, kravchenko),
            CompetitorGroupDetail(1, demidenko),
            CompetitorGroupDetail(1, novokhatskiy),
          ];
        default:
          throw NotFoundException("Группа не найдена");
      }
    });
  }

  @override
  Future<FinalGrid> getIndividualGroupFinalGrid(int groupId) {
    defineSparring(sparringId, top, bottom) => Sparring(
      sparringId,
      SparringPlace(top.id, top.shrink(), _rangeGroup, false, null, 0),
      bottom != null
          ? SparringPlace(
            bottom.id,
            bottom.shrink(),
            _rangeGroup,
            false,
            null,
            0,
          )
          : null,
      SparringState.ongoing,
    );
    var sparringList = <Sparring>[
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
          Quarterfinal(
            sparringList[0],
            sparringList[1],
            sparringList[2],
            sparringList[3],
          ),
          Semifinal(sparringList[4], sparringList[5]),
          Final(sparringList[6], sparringList[7]),
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
    return Future.delayed(delay, () {
      if (sectionId == 1) {
        return _section;
      } else {
        throw NotFoundException("Секция не найдена");
      }
    });
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
    return Future.delayed(delay, () => _rangeGroup);
  }

  @override
  Future<SparingPlace> getSparringPlace(int placeId) {
    // TODO: implement getSparringPlace
    throw UnimplementedError();
  }

  @override
  Future<RangeGroup> getSparringPlacesRanges(int placeId) {
    return Future.delayed(delay, () => _rangeGroup);
  }

  @override
  Future<int> login(Credentials credentials) async {
    if (credentials.login == "Недотёпа") {
      throw InvalidParametersException("Неверные параметры входа");
    }
    return Future.delayed(delay, () => 1);
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
    var range = _rangeGroup.ranges[request.rangeOrdinal - 1];
    range.shots = request.shots;
    return Future.delayed(delay, () => range);
  }

  @override
  Future<Range> putSparringPlacesRange(int placeId, ChangeRange request) {
    var range = _rangeGroup.ranges[request.rangeOrdinal - 1];
    range.shots = request.shots;
    return Future.delayed(delay, () => range);
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
