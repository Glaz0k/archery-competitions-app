import 'package:flutter/material.dart';

class CompetitionField extends StatelessWidget {
  final String nameOfComp;
  final String date;
  const CompetitionField({super.key, required this.nameOfComp, required this.date});

  @override
  Widget build(BuildContext context) {
    return Card(
      margin: EdgeInsets.symmetric(vertical: 10.0, horizontal: 25.0),
      child: ListTile(
        title: Text(
          nameOfComp,
          style: TextStyle(fontSize: 17, fontWeight: FontWeight.bold),
        ),
        subtitle: Text('Даты проведения $date'),
        leading: IconButton(
          icon: Icon(Icons.info_outline),
          color: Theme.of(context).primaryColor,
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