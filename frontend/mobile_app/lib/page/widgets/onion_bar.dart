import 'package:flutter/material.dart';
import 'package:mobile_app/page/profile_page.dart';

class OnionBar extends AppBar {
  OnionBar(String title, BuildContext context, {super.key, super.bottom})
      : super(
    title: Text(title),
    backgroundColor: Theme.of(context).colorScheme.inversePrimary,
    actions: [_ProfileButton()]
  );

  OnionBar.withoutProfile(String title, BuildContext context,
      {super.key, super.bottom}) : super(
      title: Text(title),
      backgroundColor: Theme.of(context).colorScheme.inversePrimary
  );
}

class _ProfileButton extends StatelessWidget{
  @override
  Widget build(BuildContext context) {
    return IconButton(onPressed: () {
      Navigator.push(context, MaterialPageRoute(builder: (context) => ProfilePage()));
    }, icon: Icon(Icons.account_circle));
  }
}