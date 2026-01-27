import { Controller, Get, Post, Body } from '@nestjs/common';
import { ApiTags, ApiOperation } from '@nestjs/swagger';
import { ContactService } from './contact.service';
import { CreateConsultationDto } from './dto/create-consultation.dto';

@ApiTags('Contact')
@Controller('contact')
export class ContactController {
  constructor(private readonly service: ContactService) {}

  @Post('consultation')
  @ApiOperation({ summary: 'Submit consultation request' })
  create(@Body() dto: CreateConsultationDto) {
    return this.service.createConsultation(dto);
  }

  @Get('consultation')
  @ApiOperation({ summary: 'Get all consultation requests' })
  findAll() {
    return this.service.findAll();
  }
}
