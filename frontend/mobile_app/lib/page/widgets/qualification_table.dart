import 'package:flutter/material.dart';
import 'package:mobile_app/api/responses.dart';

class MyQualificationTable extends StatelessWidget {
  final Section? section;

  const MyQualificationTable({super.key, required this.section});

  @override
  Widget build(BuildContext context) {
    if (section == null) {
      return const Center(child: CircularProgressIndicator());
    }
    return Container(
      decoration: BoxDecoration(
        border: Border.all(color: Theme.of(context).dividerColor),
        borderRadius: BorderRadius.circular(8),
      ),
      child: DataTable(
        columnSpacing: 50,
        columns: const [
          DataColumn(
            label: FittedBox(
              fit: BoxFit.scaleDown,
              alignment: Alignment.center,
              child: Text(
                "Место",
                style: TextStyle(fontWeight: FontWeight.bold),
              ),
            ),
            numeric: true,
          ),
          DataColumn(
            label: FittedBox(
              fit: BoxFit.scaleDown,
              alignment: Alignment.center,
              child: Text(
                "Участник",
                style: TextStyle(fontWeight: FontWeight.bold),
              ),
            ),
          ),
          DataColumn(
            label: FittedBox(
              fit: BoxFit.scaleDown,
              alignment: Alignment.center,
              child: Text(
                "Итог",
                style: TextStyle(fontWeight: FontWeight.bold),
              ),
            ),
            numeric: true,
          ),
        ],
        rows: [
          DataRow(
            cells: [
              DataCell(
                Align(
                  alignment: Alignment.center,
                  child: Text(section!.place.toString()),
                ),
              ),
              DataCell(
                FittedBox(
                  fit: BoxFit.scaleDown,
                  alignment: Alignment.center,
                  child: Text(
                    section!.competitor.fullName,
                    overflow: TextOverflow.ellipsis,
                  ),
                ),
              ),
              DataCell(
                Align(
                  alignment: Alignment.center,
                  child: Text(section!.total.toString()),
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }
}
