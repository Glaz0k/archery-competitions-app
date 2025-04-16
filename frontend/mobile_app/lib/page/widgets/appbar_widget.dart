import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';

AppBar buildAppBar(BuildContext context) {
  final moonIcon = CupertinoIcons.settings;
  return AppBar(
    title: Text(
      "Личный кабинет",
      style: TextStyle(fontSize: 13, fontWeight: FontWeight.bold),
    ),
    leading: BackButton(),
    centerTitle: true,
    backgroundColor: Colors.transparent,
    elevation: 0,
    actions: [
      IconButton(
        onPressed: () {
          Navigator.pushNamed(context, '/competition_page');
        },
        icon: Icon(moonIcon),
      ),
    ],
  );
}
