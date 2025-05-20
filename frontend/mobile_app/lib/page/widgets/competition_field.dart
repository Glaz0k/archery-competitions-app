import 'package:flutter/material.dart';
import 'package:mobile_app/page/hub_page.dart';

class CompetitionField extends StatelessWidget {
  final String nameOfComp;
  final String date;
  final int groupId;

  const CompetitionField({
    super.key,
    required this.nameOfComp,
    required this.date,
    required this.groupId,
  });

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
        trailing: IconButton(
          onPressed: () {
            Navigator.push(context, MaterialPageRoute(builder: (context) => HubPage(individualGroupId: groupId, title: "Группа $groupId")));
          },
          icon: Icon(Icons.arrow_forward),
        ),
      ),
    );
  }
}
