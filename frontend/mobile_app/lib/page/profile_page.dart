import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import 'package:keyboard_dismisser/keyboard_dismisser.dart';
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
  String formatDate(DateTime date) {
    try {
      final format = DateFormat('dd.MM.yyyy', 'ru_RU');
      return format.format(date);
    } catch (e) {
      return '${date.day.toString().padLeft(2, '0')}.${date.month.toString().padLeft(2, '0')}.${date.year}';
    }
  }

  @override
  Widget build(BuildContext context) {
    final user = Provider.of<UserProvider>(context).userPref;
    return KeyboardDismisser(
      gestures: [GestureType.onTap],
      child: GestureDetector(
        child: Scaffold(
          appBar: OnionBar("Профиль", context),
          body: SingleChildScrollView(
            physics: AlwaysScrollableScrollPhysics(), //BouncingScrollPhysics
            child: Container(
              padding: const EdgeInsets.all(8),
              child: Column(
                children: [
                  Card(
                    margin: EdgeInsets.symmetric(
                      vertical: 10.0,
                      horizontal: 25.0,
                    ),
                    child: ListTile(
                      title: const Text("Фамилия Имя"),
                      subtitle: Text(user.fullName),
                      leading: Icon(
                        CupertinoIcons.person_alt,
                        color: Colors.teal,
                      ),
                    ),
                  ),

                  Card(
                    margin: EdgeInsets.symmetric(
                      vertical: 10.0,
                      horizontal: 25.0,
                    ),
                    child: ListTile(
                      title: const Text("Дата рождения"),
                      subtitle: Text(formatDate(user.dateOfBirth)),
                      leading: Icon(Icons.date_range, color: Colors.teal),
                    ),
                  ),

                  Card(
                    margin: EdgeInsets.symmetric(
                      vertical: 10.0,
                      horizontal: 25.0,
                    ),
                    child: ListTile(
                      title: const Text("Пол"),
                      subtitle: Text(user.identity.getGender),
                      leading: Icon(Icons.perm_identity, color: Colors.teal),
                    ),
                  ),
                  Card(
                    margin: EdgeInsets.symmetric(
                      vertical: 10.0,
                      horizontal: 25.0,
                    ),
                    child: ListTile(
                      title: const Text("Лук"),
                      subtitle: Text(user.bow!.getBowClass),
                      leading: Icon(Icons.dangerous, color: Colors.teal),
                    ),
                  ),
                  Card(
                    margin: EdgeInsets.symmetric(
                      vertical: 10.0,
                      horizontal: 25.0,
                    ),
                    child: ListTile(
                      title: const Text("Спортивный разряд"),
                      subtitle: Text(user.rank!.getSportsRank),
                      leading: Icon(Icons.star, color: Colors.teal),
                    ),
                  ),
                  Card(
                    margin: EdgeInsets.symmetric(
                      vertical: 10.0,
                      horizontal: 25.0,
                    ),
                    child: ListTile(
                      title: const Text("Город"),
                      subtitle: Text(user.region),
                      leading: Icon(Icons.location_city, color: Colors.teal),
                    ),
                  ),
                  Card(
                    margin: EdgeInsets.symmetric(
                      vertical: 10.0,
                      horizontal: 25.0,
                    ),
                    child: ListTile(
                      title: const Text("Клуб"),
                      subtitle: Text(user.club),
                      leading: Icon(
                        CupertinoIcons.sportscourt_fill,
                        color: Colors.teal,
                      ),
                    ),
                  ),
                  ProfileWidget(onClicked: () {}),
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
                        // Navigator.pushNamed(
                        //   context,
                        //   '/edit_profile_page',
                        //   arguments: user,
                        // );
                      },
                      // style: ElevatedButton.styleFrom(
                      //   backgroundColor: Colors.green,
                      //   side: BorderSide.none,
                      //   shape: StadiumBorder(),
                      //   minimumSize: Size(220, 50),
                      // ),
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
        ),
      ),
    );
  }
}

class ProfileWidget extends StatelessWidget {
  final VoidCallback onClicked;

  const ProfileWidget({super.key, required this.onClicked});

  @override
  Widget build(BuildContext context) {
    final color = Theme.of(context).colorScheme.primary;
    return Center(child: Stack(children: []));
  }

  Widget buildEditIcon(Color color) => buildCircle(
    color: Colors.white,
    all: 3,
    child: buildCircle(
      color: color,
      all: 8,
      child: Icon(Icons.edit, size: 20, color: Colors.white),
    ),
  );

  Widget buildCircle({
    required Widget child,
    required double all,
    required Color color,
  }) => ClipOval(
    child: Container(color: color, padding: EdgeInsets.all(all), child: child),
  );
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
