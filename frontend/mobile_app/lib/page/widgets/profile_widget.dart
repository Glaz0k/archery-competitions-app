import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';

class ProfileWidget extends StatelessWidget {
  final VoidCallback onClicked;

  const ProfileWidget({Key? key, required this.onClicked}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final color = Theme.of(context).colorScheme.primary;
    return Center(
      child: Stack(
        children: [
        ],
      ),
    );
  }
   Widget buildEditIcon(Color color) => buildCircle(
       color: Colors.white,
        all: 3,
        child: buildCircle(
          color: color,
          all: 8,
          child: Icon(
            Icons.edit,
            size: 20,
            color: Colors.white,
          ),
        ),
   );


  Widget buildCircle ({required Widget child, required double all, required Color color}) =>
      ClipOval(
        child: Container(
            color: color,
            padding: EdgeInsets.all(all),
            child: child,
      )
  );
}

bool isValidPhoneNumber(String num) {
  return RegExp(r'^\\+?\\d[ ]?[-(]?\\d{3}[-)]?[ ]?\\d{3}[- ]?\\d{2}[- ]?\\d{2}$').hasMatch(num);
}

bool isValidEmail (String email) {
  return RegExp(r'^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$').hasMatch(email);
}