import 'package:flutter/material.dart';
import 'package:mobile_app/page/history_page.dart';
import 'package:mobile_app/page/profile_page.dart';
import 'package:mobile_app/page/series_input_page.dart';

void main() => runApp(
  MaterialApp(
    theme: ThemeData(
      colorScheme: ColorScheme.fromSeed(seedColor: Colors.deepPurple),
    ),
    home: Onion(),
  ),
);

class Onion extends StatefulWidget {
  const Onion({super.key});

  @override
  State<StatefulWidget> createState() => _OnionState();
}

class _OnionState extends State<Onion> {
  static const List<NavigationDestination> _destinations = [
    NavigationDestination(
      selectedIcon: Icon(Icons.history_edu),
      icon: Icon(Icons.history_edu_outlined),
      label: 'История',
    ),
    NavigationDestination(
      selectedIcon: Icon(Icons.ads_click),
      icon: Icon(Icons.ads_click_outlined),
      label: 'Результаты',
    ),
    NavigationDestination(
      selectedIcon: Icon(Icons.account_circle),
      icon: Icon(Icons.account_circle_outlined),
      label: 'Профиль',
    ),
  ];
  late final List<Widget> _mainPages;
  int _currentPage = 1;

  @override
  void initState() {
    super.initState();
    _mainPages = [HistoryPage(), SeriesInputPage(), ProfilePage()];
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: _mainPages[_currentPage],
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
