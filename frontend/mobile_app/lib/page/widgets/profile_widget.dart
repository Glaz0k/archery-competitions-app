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