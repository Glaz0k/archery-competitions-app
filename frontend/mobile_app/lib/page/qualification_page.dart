import 'package:flutter/material.dart';
import 'package:mobile_app/page/widgets/onion_bar.dart';
import 'package:mobile_app/page/widgets/qualification_table.dart';
import 'package:provider/provider.dart';

import '../api/api.dart';
import '../api/responses.dart';

class QualificationPage extends StatefulWidget {
  final int sectionId;

  const QualificationPage({super.key, required this.sectionId});

  @override
  State<QualificationPage> createState() => _QualificationPageState();
}

class _QualificationPageState extends State<QualificationPage> {
  Section? _section;

  @override
  void initState() {
    super.initState();
    _loadData();
  }

  Future<void> _loadData() async {
    var api = context.read<Api>();
    final section = await api.getQualificationSection(widget.sectionId);
    setState(() {
      _section = section;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: OnionBar("Квалификация", context),
      body: RefreshIndicator(
        onRefresh: _loadData,
        child: SingleChildScrollView(
          physics: AlwaysScrollableScrollPhysics(),
          child: Padding(
            padding: const EdgeInsets.all(16.0),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  "Секция ${_section?.id}",
                  style: Theme.of(context).textTheme.titleLarge,
                ),
                const SizedBox(height: 16),
                MyQualificationTable(section: _section),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
