import 'package:flutter/material.dart';

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
   String fullName;
   DateTime dateOfBirth;
   Gender identity;
   BowClass bow;
   SportsRank rank;
   String region;
   String federation;
   String club;
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
    dateOfBirth: DateTime.now(),
    identity: Gender.male,
    bow: BowClass.classic,
    rank: SportsRank.master,
    region: "SPB",
    federation: "Example",
    club: "Polytech",
  );
}

class UserProvider with ChangeNotifier {
  var _user = UserPreferences.myUser;

  User get userPref => _user;

  void updateUser(User newUser) {
    _user = newUser;
    notifyListeners();
  }

  void updateFullName(String newName) {
    _user.fullName = newName;
    notifyListeners();
  }

  void updateDateOfBirth(DateTime newDate) {
    _user.dateOfBirth = newDate;
    notifyListeners();
  }

  void updateGender(Gender newGender) {
    _user.identity = newGender;
    notifyListeners();
  }

  void updateBow(BowClass newBow) {
    _user.bow = newBow;
    notifyListeners();
  }

  void updateRank(SportsRank newRank) {
    _user.rank = newRank;
    notifyListeners();
  }

  void updateRegion(String newRegion) {
    _user.region = newRegion;
    notifyListeners();
  }

  void updateClub(String newClub) {
    _user.club = newClub;
    notifyListeners();
  }
}

extension GenderExtension on Gender {
  String get getGender {
    switch(this) {
      case Gender.male: return "Мужчина";
      case Gender.female: return "Женщина";
    }
  }
}

extension BowExtension on BowClass {
  String get getBowClass {
    switch(this) {
      case BowClass.classic: return "Классический";
      case BowClass.block: return "Блочный";
      case BowClass.classicNewbie: return "";
      case BowClass.classic3D: return "";
      case BowClass.compound3D: return "";
      case BowClass.long3D: return "";
      case BowClass.peripheral: return "";
      case BowClass.peripheralWithRing: return "";
    }
  }
}

extension SportsRankExtension on SportsRank {
  String get getSportsRank {
    switch(this) {
      case SportsRank.master: return "Мастер спорта";
      case SportsRank.candidateForMaster: return "Кандидат мастера спорта";
      case SportsRank.firstClass: return "Первый разряд";
      case SportsRank.secondClass: return "Второй разряд";
      case SportsRank.thirdClass: return "Третий разряд";
      case SportsRank.masterInternational: return "Международный магистр";
      case SportsRank.meritedMaster: return "Заслуженный мастер спорта";
    }
  }
}