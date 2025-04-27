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
   BowClass? bow;
   SportsRank? rank;
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
    identity: Gender.female,
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

  void updateBow(BowClass? newBow) {
    _user.bow = newBow;
    notifyListeners();
  }

  void updateRank(SportsRank? newRank) {
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
      case BowClass.classicNewbie: return "КЛ(новички)";
      case BowClass.classic3D: return "3Д-классический лук";
      case BowClass.compound3D: return "3Д-составной лук";
      case BowClass.long3D: return "3Д-длинный лук";
      case BowClass.peripheral: return "Периферийный лук";
      case BowClass.peripheralWithRing: return "Периферийный лук(с кольцом)";
    }
  }
  static BowClass? setBowClass(String? value) {
    switch(value) {
      case "Классический": return BowClass.classic;
      case "Блочный": return BowClass.block;
      case "КЛ(новички)": return BowClass.classicNewbie;
      case "3Д-классический лук": return BowClass.classic3D;
      case "3Д-составной лук": return BowClass.compound3D;
      case "3Д-длинный лук": return BowClass.long3D;
      case "Периферийный лук": return BowClass.peripheral;
      case "Периферийный лук(с кольцом)": return BowClass.peripheralWithRing;
      default: return null;
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
  static SportsRank? setSportsRank(String? val) {
    switch(val) {
      case "Заслуженный мастер спорта": return SportsRank.meritedMaster;
      case "Мастер спорта международного класса": return SportsRank.masterInternational;
      case "Мастер спорта": return SportsRank.master;
      case "Кандидат в мастера спорта": return SportsRank.candidateForMaster;
      case "Первый разряд": return SportsRank.firstClass;
      case "Второй разряд": return SportsRank.secondClass;
      case "Третий разряд": return SportsRank.thirdClass;
      default: return null;
    }
  }
}