class User {
  final String name;
  final String surname;
  final String middleName;
  final String phoneNumber;
  final String city;
  final String login;
  final String club;
  const User({required this.name, required this.surname, required this.middleName, required this.phoneNumber, required this.city, required this.login, required this.club,});
}

class UserPreferences {
  static const myUser = User(name: "Danil",
      surname: "Novo",
      middleName: "God",
      phoneNumber: "321421",
      city: "SPB",
      login: "login",
      club: "Polytech",
  );
}