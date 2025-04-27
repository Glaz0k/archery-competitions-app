import 'package:flutter/material.dart';

class OnionBar extends AppBar {
  OnionBar(String title, BuildContext context, {super.key}) : super(
    title: Text(title),
    backgroundColor: Theme.of(context).colorScheme.inversePrimary
  );
}