import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';

AppBar buildAppBar(BuildContext context) {
  final moonIcon = CupertinoIcons.moon_stars;
  return AppBar(
    title: Text("Личный кабинет", style: TextStyle(fontSize: 13, fontWeight: FontWeight.bold),),
    leading: BackButton(),
    backgroundColor: Colors.transparent,
    elevation: 0,
    actions: [
      IconButton(onPressed: () {},
          icon: Icon(moonIcon)
      )
    ],
  );
}