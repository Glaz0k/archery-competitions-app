import 'package:flutter/cupertino.dart';

class MyTextBox extends StatelessWidget {
  final String text;
  final String sectionName;
  const MyTextBox({super.key, required this.text, required this.sectionName});

  @override
  Widget build(BuildContext context) {
    return Column(children: [Text(sectionName), Text(text)]);
  }
}
