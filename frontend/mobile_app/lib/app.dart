import 'package:flutter/material.dart';
import 'package:mobile_app/page/login_page.dart';
import 'package:mobile_app/page/main_competition_page.dart';
import 'package:mobile_app/page/widgets/user.dart';
import 'package:provider/provider.dart';

import 'api/api.dart';

class OnionApp extends StatelessWidget {
  const OnionApp({super.key});

  @override
  Widget build(BuildContext context) {
    var userProvider = context.watch<UserProvider>();
    var api = context.watch<Api>();
    loadingScreen() => Center(child: CircularProgressIndicator());
    if (userProvider.loading) {
      return loadingScreen();
    } else {
      if (userProvider.getId() == null) {
        return LoginPage();
      } else {
        return FutureBuilder(
          future: userProvider.loadUser(api),
          builder: (context, snapshot) {
            if (snapshot.connectionState == ConnectionState.done) {
              return MainCompetitionPage();
            } else {
              return loadingScreen();
            }
          },
        );
      }
    }
  }
}
