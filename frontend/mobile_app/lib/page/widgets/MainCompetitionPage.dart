import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:keyboard_dismisser/keyboard_dismisser.dart';

class MainCompetitionPage extends StatefulWidget {
  const MainCompetitionPage({super.key});

  @override
  State<MainCompetitionPage> createState() => _MainCompetitionPage();
}

class _MainCompetitionPage extends State<MainCompetitionPage> {
  @override
  Widget build(BuildContext context) {
    return KeyboardDismisser(
      gestures: [GestureType.onTap],
      child: GestureDetector(
        child: Scaffold(
          appBar: AppBar(
            title: Text(
              "Кубки",
              style: TextStyle(fontSize: 15, color: Colors.black87),
            ),
            leading: BackButton(),
            centerTitle: true,
            backgroundColor: Colors.green,
          ),
          body: Center(child: ListView(children: [
            
            ],
            )),
        ),
      ),
    );
  }

  Widget buildCompitionField(String nameOfComp, String address, String season) {
    return Card(
      margin: EdgeInsets.symmetric(vertical: 10.0, horizontal: 25.0),
      child: ListTile(
        title: Text(nameOfComp),
        subtitle: Text("Адрес: $address\nСезон: $season"),
        leading: IconButton(
          icon: Icon(Icons.info_outline),
          color: Colors.teal,
          onPressed: () {},
        ),
        trailing: IconButton(
          onPressed: () {
            Navigator.pushNamed(context, "/profile_page");
          },
          icon: Icon(Icons.join_full),
        ),
      ),
    );
  }
}
