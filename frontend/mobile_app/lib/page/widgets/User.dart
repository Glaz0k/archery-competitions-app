class User {
  final String name;
  final String surname;
  final String middleName;
  final String phoneNumber;
  final String city;
  final String email;
  final String club;
  final bool isSettings;
  const User({required this.name, required this.surname, required this.middleName, required this.phoneNumber, required this.city, required this.email, required this.club, required this.isSettings});
}

class UserPreferences {
  static const myUser = User(name: "Danil",
      surname: "Novo",
      middleName: "God",
      phoneNumber: "321421",
      city: "SPB",
      email: "faskl@mail.ru",
      club: "Polytech",
      isSettings: false
  );
}