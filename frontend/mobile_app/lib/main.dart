import 'package:flutter/material.dart';
import 'package:mobile_app/page/competition_page.dart';
import 'package:mobile_app/page/profile_page.dart';
import 'package:mobile_app/page/series_page.dart';
import 'package:provider/provider.dart';

import 'model/series_model.dart';

void main() => runApp(
  MaterialApp(
    theme: ThemeData(
      colorScheme: ColorScheme.fromSeed(
        seedColor: Colors.blue,
        secondary: Colors.red,
      ),
    ),
    home: MultiProvider(
      providers: [ChangeNotifierProvider(create: (context) => SeriesModel())],
      child: Onion(),
    ),
  ),
);

class Onion extends StatefulWidget {
  Onion({super.key});

  final List<Widget> _mainPages = [
    CompetitionPage(),
    SeriesPage(),
    ProfilePage(),
  ];

  @override
  State<StatefulWidget> createState() => _OnionState();
}

class _OnionState extends State<Onion> {
  static const List<NavigationDestination> _destinations = [
    NavigationDestination(
      selectedIcon: Icon(Icons.info),
      icon: Icon(Icons.info_outline),
      label: 'Объявления',
    ),
    NavigationDestination(
      selectedIcon: Icon(Icons.scoreboard),
      icon: Icon(Icons.scoreboard_outlined),
      label: 'Серии',
    ),
    NavigationDestination(
      selectedIcon: Icon(Icons.account_circle),
      icon: Icon(Icons.account_circle_outlined),
      label: 'Профиль',
    ),
  ];

  int _currentPage = 1;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: widget._mainPages[_currentPage],
      bottomNavigationBar: NavigationBar(
        selectedIndex: _currentPage,
        destinations: _destinations,
        onDestinationSelected:
            (idx) => setState(() {
              _currentPage = idx;
            }),
        //labelBehavior: NavigationDestinationLabelBehavior.alwaysHide,
      ),
    );
  }
}
