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
    if (userProvider.loading) {
      return SizedBox.shrink(); // Пустой экран
    } else {
      if (userProvider.getUser(api) == null) {
        return LoginPage();
      } else {
        return MainCompetitionPage();
      }
    }
  }
}
