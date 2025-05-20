import 'package:flutter/material.dart';
import 'package:mobile_app/page/edit_profile_page.dart';
import 'package:mobile_app/page/widgets/onion_bar.dart';
import 'package:provider/provider.dart';

import 'widgets/user.dart';

class ProfilePage extends StatefulWidget {
  const ProfilePage({super.key});

  @override
  State<ProfilePage> createState() => _ProfilePageState();
}

class _ProfilePageState extends State<ProfilePage> {
  @override
  Widget build(BuildContext context) {
    var user = context.watch<UserProvider>().user;
    return Scaffold(
      appBar: OnionBar.withoutProfile("Профиль", context),
      body: SingleChildScrollView(
        physics: AlwaysScrollableScrollPhysics(),
        child: Container(
          padding: const EdgeInsets.all(8),
          child: Column(
            children: [
              Card(
                margin: EdgeInsets.symmetric(vertical: 10.0, horizontal: 25.0),
                child: ListTile(
                  title: const Text(
                    "Фамилия Имя",
                    style: TextStyle(fontWeight: FontWeight.bold),
                  ),
                  subtitle: Text(user?.fullName ?? ""),
                  leading: Icon(
                    Icons.person,
                    color: Theme.of(context).primaryColor,
                  ),
                ),
              ),

              Card(
                margin: EdgeInsets.symmetric(vertical: 10.0, horizontal: 25.0),
                child: ListTile(
                  title: const Text(
                    "Дата рождения",
                    style: TextStyle(fontWeight: FontWeight.bold),
                  ),
                  subtitle: Text(user?.birthDate ?? ''),
                  leading: Icon(
                    Icons.date_range,
                    color: Theme.of(context).primaryColor,
                  ),
                ),
              ),

              Card(
                margin: EdgeInsets.symmetric(vertical: 10.0, horizontal: 25.0),
                child: ListTile(
                  title: const Text(
                    "Пол",
                    style: TextStyle(fontWeight: FontWeight.bold),
                  ),
                  subtitle: Text(user?.identity.toString() ?? ''),
                  leading: Icon(
                    Icons.perm_identity,
                    color: Theme.of(context).primaryColor,
                  ),
                ),
              ),
              Card(
                margin: EdgeInsets.symmetric(vertical: 10.0, horizontal: 25.0),
                child: ListTile(
                  title: const Text(
                    "Лук",
                    style: TextStyle(fontWeight: FontWeight.bold),
                  ),
                  subtitle: Text(user?.bow.toString() ?? ''),
                  leading: Icon(
                    Icons.subject,
                    color: Theme.of(context).primaryColor,
                  ),
                ),
              ),
              Card(
                margin: EdgeInsets.symmetric(vertical: 10.0, horizontal: 25.0),
                child: ListTile(
                  title: const Text(
                    "Спортивный разряд",
                    style: TextStyle(fontWeight: FontWeight.bold),
                  ),
                  subtitle: Text(user?.rank.toString() ?? ''),
                  leading: Icon(
                    Icons.star,
                    color: Theme.of(context).primaryColor,
                  ),
                ),
              ),
              Card(
                margin: EdgeInsets.symmetric(vertical: 10.0, horizontal: 25.0),
                child: ListTile(
                  title: const Text(
                    "Город",
                    style: TextStyle(fontWeight: FontWeight.bold),
                  ),
                  subtitle: Text(user?.region ?? ''),
                  leading: Icon(
                    Icons.location_city,
                    color: Theme.of(context).primaryColor,
                  ),
                ),
              ),
              Card(
                margin: EdgeInsets.symmetric(vertical: 10.0, horizontal: 25.0),
                child: ListTile(
                  title: const Text(
                    "Федерация",
                    style: TextStyle(fontWeight: FontWeight.bold),
                  ),
                  subtitle: Text(user?.federation ?? ''),
                  leading: Icon(
                    Icons.people_alt,
                    color: Theme.of(context).primaryColor,
                  ),
                ),
              ),
              Card(
                margin: EdgeInsets.symmetric(vertical: 10.0, horizontal: 25.0),
                child: ListTile(
                  title: const Text(
                    "Клуб",
                    style: TextStyle(fontWeight: FontWeight.bold),
                  ),
                  subtitle: Text(user?.club ?? ''),
                  leading: Icon(
                   Icons.scoreboard,
                    color: Theme.of(context).primaryColor,
                  ),
                ),
              ),
              SizedBox(height: 30),
              SizedBox(
                child: FilledButton(
                  onPressed: () {
                    Navigator.push(
                      context,
                      MaterialPageRoute(
                        builder: (context) => EditProfilePage(),
                      ),
                    );
                  },
                  child: const Text(
                    "Редактировать",
                    style: TextStyle(color: Colors.white),
                  ),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

bool isValidPhoneNumber(String num) {
  return RegExp(
    r'^\\+?\\d?[-(]?\\d{3}[-)]??\\d{3}[- ]?\\d{2}[- ]?\\d{2}$',
  ).hasMatch(num);
}

bool isValidEmail(String email) {
  return RegExp(
    r'^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$',
  ).hasMatch(email);
}
