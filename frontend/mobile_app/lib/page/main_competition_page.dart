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
              "Соревнования",
              style: Theme.of(context).appBarTheme.titleTextStyle
            ),
            leading: BackButton(),
            centerTitle: true,
            backgroundColor: Colors.green,
          ),
          body: Center(child: ListView(children: [
            buildCompitionField("| этап", "21-22 октября 2023"),
            buildCompitionField("|| этап", "21-22 октября 2024"),
            buildCompitionField("||| этап", "21-22 октября 2025"),
            buildCompitionField("Финал", "21-22 ноябре 2025"),
            ],
            )),
        ),
      ),
    );
  }

  Widget buildCompitionField(String nameOfComp, String date) {
    return Card(
      margin: EdgeInsets.symmetric(vertical: 10.0, horizontal: 25.0),
      child: ListTile(
        title: Text(nameOfComp, style: TextStyle(fontSize: 17, fontWeight: FontWeight.bold),),
        subtitle: Text('Даты проведения: $date г.'),
        leading: IconButton(
          icon: Icon(Icons.info_outline),
          color: Colors.teal,
          onPressed: () {
            Navigator.pushNamed(context, '/individual_group');
          },
        ),
        trailing: IconButton(
          onPressed: () {
            Navigator.pushNamed(context, "/profile_page");
          },
          icon: Icon(Icons.person),
        ),
      ),
    );
  }
}
