import 'package:flutter/material.dart';
import 'package:keyboard_dismisser/keyboard_dismisser.dart';
import 'package:mobile_app/api/responses.dart';
import 'package:mobile_app/page/widgets/onion_bar.dart';
import 'package:mobile_app/api/real.dart';

class MainCompetitionPage extends StatefulWidget {
  const MainCompetitionPage({super.key});

  @override
  State<MainCompetitionPage> createState() => _MainCompetitionPage();
}

class _MainCompetitionPage extends State<MainCompetitionPage> {

  RealServer api = RealServer();
  late Future<List<Competition>> competitionsFuture = _loadCompetitions();

  @override
  Widget build(BuildContext context) {
    return KeyboardDismisser(
      gestures: [GestureType.onTap],
      child: GestureDetector(
        child: Scaffold(
          appBar: OnionBar("Соревнования", context),
          body: Center(
            child: FutureBuilder(future: competitionsFuture,
                builder: (context, snapshot) {
                  if (snapshot.connectionState == ConnectionState.waiting) {
                    return CircularProgressIndicator();
                  } else if (snapshot.hasError) {
                    return Text("Ошибка: ${snapshot.error}");
                  } else if (!snapshot.hasData || snapshot.data!.isEmpty) {
                    return const Text("Нет доступных соревнований");
                  } else {
                    final competitions = snapshot.data!;
                    return ListView.builder(itemCount: competitions.length,
                      itemBuilder: (context, index) {
                        final competition = competitions[index];
                        return buildCompetitionField(competition.stage.name, _formatCompetitionDate(competition.startDate, competition.endDate));
                      },
                    );
                  }
                }
            ),
          ),
        ),
      ),
    );
  }

  Widget buildCompetitionField(String nameOfComp, String date) {
    return Card(
      margin: EdgeInsets.symmetric(vertical: 10.0, horizontal: 25.0),
      child: ListTile(
        title: Text(
          nameOfComp,
          style: TextStyle(fontSize: 17, fontWeight: FontWeight.bold),
        ),
        subtitle: Text('Даты проведения: $date'),
        leading: IconButton(
          icon: Icon(Icons.info_outline),
          color: Colors.teal,
          onPressed: () {
            Navigator.pushNamed(context, '/individual_group');
          },
        ),
        trailing: IconButton(
          onPressed: () {
            Navigator.pushNamed(context, "/profile_page");
          },
          icon: Icon(Icons.person),
        ),
      ),
    );
  }

  Future<List<Competition>> _loadCompetitions() async {
    try {
      List<Cup> cups = await api.getCups();
      if (cups.isEmpty) return [];
      int cupId = cups.first.id;
      return await api.getCupsCompetitions(cupId);
    } catch (e) {
      throw Exception('Не удалось загрузить соревнования');
    }
  }

  String _formatCompetitionDate(DateTime? start, DateTime? end) {
    if (start == null && end == null) {
      return "Даты проведения соревнований не указаны";
    }
    if (start == null) return 'до ${_formatDate(end!)}';
    if (end == null) return 'с ${_formatDate(start)}';

    if (start.day == end.day && start.month == end.month &&
        start.year == end.year) {
      return _formatDate(start);
    }

    if (start.year == end.year && start.month == end.month) {
      return '${start.day}-${end.day} ${_getMonth(start.month)} ${start
          .year} г.';
    }

    if (start.year == end.year) {
      return '${start.day} ${_getMonth(start.month)} - '
          '${end.day} ${_getMonth(end.month)} ${start.year} г.';
    }

    return '${_formatDate(start)} - ${_formatDate(end)}';
  }

  String _getMonth(int month) {
    const months = ['января', 'февраля', 'марта', 'апреля', 'мая', 'июня',
      'июля', 'августа', 'сентября', 'октября', 'ноября', 'декабря'];
    return months[month - 1];
  }

  String _formatDate(DateTime date) {
    return '${date.day} ${_getMonth(date.month)} ${date.year} г.';
  }
}
