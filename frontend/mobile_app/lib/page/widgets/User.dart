enum SportsRank {
  meritedMaster,
  masterInternational,
  master,
  candidateForMaster,
  firstClass,
  secondClass,
  thirdClass,
}

enum Gender { male, female }

enum BowClass {
  classic,
  block,
  classicNewbie,
  classic3D,
  compound3D,
  long3D,
  peripheral,
  peripheralWithRing,
}

class User {
  final int id;
  final String fullName;
  final DateTime dateOfBirth;
  final Gender identity;
  final BowClass bow;
  final SportsRank rank;
  final String region;
  final String? federation;
  final String? club;
  User({
    required this.id,
    required this.fullName,
    required this.dateOfBirth,
    required this.identity,
    required this.bow,
    required this.rank,
    required this.region,
    required this.federation,
    required this.club,
  });
}

class UserPreferences {
  static var myUser = User(
    id: 1,
    fullName: "Novo danil",
    dateOfBirth: DateTime.parse("09.02.2005"),
    identity: Gender.male,
    bow: BowClass.classic,
    rank: SportsRank.master,
    region: "SPB",
    federation: null,
    club: null,
  );
}
