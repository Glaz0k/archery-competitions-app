import 'package:flutter/material.dart';
import 'package:mobile_app/page/widgets/EditProfilePage.dart';
import 'package:mobile_app/page/widgets/IndividualGroup.dart';
import 'package:mobile_app/page/widgets/LoginPage.dart';
import 'package:mobile_app/page/widgets/User.dart';
import 'package:mobile_app/page/widgets/ProfilePage.dart';
import 'package:mobile_app/page/widgets/MainCompetitionPage.dart';
import 'package:provider/provider.dart';

void main() {
  runApp(const AuthPage());
}

class AuthPage extends StatelessWidget {
  const AuthPage({super.key});

  @override
  Widget build(BuildContext context) {
    return ChangeNotifierProvider(
      create: (context) => UserProvider(),
      child: MaterialApp(
        debugShowCheckedModeBanner: false,
        initialRoute: '/',
        routes: {
          '/': (context) => const LoginPage(),
          '/profile_page': (context) => ProfilePage(),
          '/edit_profile_page': (context) => EditProfilePage(),
          '/competition_page': (context) => MainCompetitionPage(),
          '/individual_group': (context) => IndividualGroup(),
        },
        theme: ThemeData(
          //colorSchemeSeed: Colors.green,
          textTheme: TextTheme(
            headlineMedium: TextStyle(fontSize: 17, color: Colors.black)
          ),
          appBarTheme: AppBarTheme(
            titleTextStyle: TextStyle(fontSize: 18, fontWeight: FontWeight.w900, color: Colors.white),
            backgroundColor: Colors.green,
            shadowColor: Colors.green[200],
            centerTitle: true,
          )
        ),
      ),
    );
  }
}
